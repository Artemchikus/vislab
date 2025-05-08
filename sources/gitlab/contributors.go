package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vislab/sources/gitlab/types"
)

type ContributorService struct {
	client *Client
}

func (s *ContributorService) List(ctx context.Context, options *types.ListContributorsOptions, projectId int64) ([]*types.Contributor, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/contributors", projectId), options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.Contributor
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *ContributorService) ListAll(ctx context.Context, options *types.ListContributorsOptions, projectId int64) ([]*types.Contributor, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allContributors []*types.Contributor
	for {
		select {
		case <-ticker.C:
			contributors, resp, err := s.List(ctx, options, projectId)
			if err != nil {
				return nil, resp, err
			}

			allContributors = append(allContributors, contributors...)

			if resp.NextPage == 0 {
				return allContributors, nil, nil
			}
			options.Page = int64(resp.NextPage)
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}
