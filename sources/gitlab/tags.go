package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vislab/libs/ptr"
	"vislab/sources/gitlab/types"
)

type TagService struct {
	client *Client
}

func (s *TagService) List(ctx context.Context, options *types.ListTagsOptions, projectId int64) ([]*types.Tag, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/tags", projectId), options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.Tag
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *TagService) ListAll(ctx context.Context, options *types.ListTagsOptions, projectId int64) ([]*types.Tag, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allTags []*types.Tag
	for {
		select {
		case <-ticker.C:
			tags, resp, err := s.List(ctx, options, projectId)
			if err != nil {
				return nil, resp, err
			}

			allTags = append(allTags, tags...)

			if resp.NextPage == 0 {
				return allTags, nil, nil
			}
			options.Page = int64(resp.NextPage)
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}

func (s *TagService) Get(ctx context.Context, name string, projectId int64) (*types.Tag, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/tags/%s", projectId, name), nil)
	if err != nil {
		return nil, nil, err
	}

	var p *types.Tag
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *TagService) GetLatest(ctx context.Context, projectId int64) (*types.Tag, *Response, error) {
	p, resp, err := s.List(ctx, &types.ListTagsOptions{
		ListOptions: types.ListOptions{
			PerPage: 1,
		},
		OrderBy: ptr.Ptr("updated"),
		Sort:    ptr.Ptr("desc"),
	}, projectId)
	if err != nil {
		return nil, resp, err
	}

	if len(p) == 0 {
		return nil, resp, fmt.Errorf("tags not found")
	}

	return p[0], resp, nil
}

func (s *TagService) Compare(ctx context.Context, from, to string, projectId int64) (*types.CompareResult, *Response, error) {
	options := &types.CompareOptions{
		From: &from,
		To:   &to,
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/compare", projectId), options)
	if err != nil {
		return nil, nil, err
	}

	var p *types.CompareResult
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}
