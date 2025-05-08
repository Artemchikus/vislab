package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vislab/sources/gitlab/types"
)

type GroupService struct {
	client *Client
}

func (s *GroupService) List(ctx context.Context, options *types.ListGroupsOptions) ([]*types.Group, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, "groups", options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.Group
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *GroupService) ListAll(ctx context.Context, options *types.ListGroupsOptions) ([]*types.Group, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allGroups []*types.Group
	for {
		select {
		case <-ticker.C:
			Groups, resp, err := s.List(ctx, options)
			if err != nil {
				return nil, resp, err
			}

			allGroups = append(allGroups, Groups...)

			if resp.NextPage == 0 {
				return allGroups, resp, nil
			}
			options.Page = int64(resp.NextPage)
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}

func (s *GroupService) Get(ctx context.Context, id int64) (*types.Group, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("groups/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var p *types.Group
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *GroupService) ListProjects(ctx context.Context, id int64, options *types.ListProjectsOptions) ([]*types.Project, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("groups/%d/projects", id), options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.Project
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *GroupService) ListAllProjects(ctx context.Context, id int64, options *types.ListProjectsOptions) ([]*types.Project, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allProjects []*types.Project
	for {
		select {
		case <-ticker.C:
			projects, resp, err := s.ListProjects(ctx, id, options)
			if err != nil {
				return nil, resp, err
			}

			allProjects = append(allProjects, projects...)

			if resp.NextPage == 0 {
				return allProjects, nil, nil
			}
			options.Page = int64(resp.NextPage)
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}
