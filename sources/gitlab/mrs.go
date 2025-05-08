package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vislab/sources/gitlab/types"
)

type MergeRequestService struct {
	client *Client
}

func (s *MergeRequestService) List(ctx context.Context, options *types.ListMergeRequestsOptions, projectId int64) ([]*types.MergeRequest, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/merge_requests", projectId), options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.MergeRequest
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *MergeRequestService) ListAll(ctx context.Context, options *types.ListMergeRequestsOptions, projectId int64) ([]*types.MergeRequest, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allMergeRequests []*types.MergeRequest
	for {
		select {
		case <-ticker.C:
			mergeRequests, resp, err := s.List(ctx, options, projectId)
			if err != nil {
				return nil, resp, err
			}

			allMergeRequests = append(allMergeRequests, mergeRequests...)

			if resp.NextPage == 0 {
				return allMergeRequests, nil, nil
			}
			options.Page = int64(resp.NextPage)
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}

func (s *MergeRequestService) ListByTargetBranch(ctx context.Context, projectId int64, targetBranch string) ([]*types.MergeRequest, *Response, error) {
	options := &types.ListMergeRequestsOptions{
		TargetBranch: &targetBranch,
	}

	return s.ListAll(ctx, options, projectId)
}

func (s *MergeRequestService) ListBySourceBranch(ctx context.Context, projectId int64, sourceBranch string) ([]*types.MergeRequest, *Response, error) {
	options := &types.ListMergeRequestsOptions{
		SourceBranch: &sourceBranch,
	}

	return s.ListAll(ctx, options, projectId)
}
