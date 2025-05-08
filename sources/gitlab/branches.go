package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"vislab/sources/gitlab/types"
)

type BranchService struct {
	client *Client
}

func (s *BranchService) List(ctx context.Context, options *types.ListBranchesOptions, projectId int64) ([]*types.Branch, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/branches", projectId), options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.Branch
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *BranchService) ListAll(ctx context.Context, options *types.ListBranchesOptions, projectId int64) ([]*types.Branch, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allBranches []*types.Branch
	for {
		select {
		case <-ticker.C:
			branches, resp, err := s.List(ctx, options, projectId)
			if err != nil {
				return nil, resp, err
			}

			allBranches = append(allBranches, branches...)

			if resp.NextPage == 0 {
				return allBranches, nil, nil
			}
			options.Page = int64(resp.NextPage)
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}

func (s *BranchService) Get(ctx context.Context, name string, projectId int64) (*types.Branch, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/branches/%s", projectId, url.QueryEscape(name)), nil)
	if err != nil {
		return nil, nil, err
	}

	var p *types.Branch
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}
