package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vislab/sources/gitlab/types"
)

type CommitService struct {
	client *Client
}

func (s *CommitService) List(ctx context.Context, options *types.ListCommitsOptions, projectId int64) ([]*types.Commit, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/commits", projectId), options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.Commit
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *CommitService) ListAll(ctx context.Context, options *types.ListCommitsOptions, projectId int64) ([]*types.Commit, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allCommits []*types.Commit
	for {
		select {
		case <-ticker.C:
			commits, resp, err := s.List(ctx, options, projectId)
			if err != nil {
				return nil, resp, err
			}

			allCommits = append(allCommits, commits...)

			if resp.NextPage == 0 {
				return allCommits, nil, nil
			}
			options.Page = int64(resp.NextPage)
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}

func (s *CommitService) Get(ctx context.Context, sha string, projectId int64) (*types.Commit, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/commits/%s", projectId, sha), nil)
	if err != nil {
		return nil, nil, err
	}

	var p *types.Commit
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, nil
}

func (s *CommitService) ListRefs(ctx context.Context, sha string, projectId int64, options *types.ListCommitsOptions) ([]*types.Ref, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/commits/%s/refs", projectId, sha), options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.Ref
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, nil
}

func (s *CommitService) ListRefsAll(ctx context.Context, sha string, projectId int64, options *types.ListCommitsOptions) ([]*types.Ref, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allRefs []*types.Ref
	for {
		select {
		case <-ticker.C:
			refs, resp, err := s.ListRefs(ctx, sha, projectId, options)
			if err != nil {
				return nil, resp, err
			}

			allRefs = append(allRefs, refs...)

			if resp.NextPage == 0 {
				return allRefs, nil, nil
			}
			options.Page = int64(resp.NextPage)
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}
