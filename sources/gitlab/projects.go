package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"vislab/sources/gitlab/types"
)

type ProjectService struct {
	client *Client
}

func (s *ProjectService) List(ctx context.Context, options *types.ListProjectsOptions) ([]*types.Project, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, "projects", options)
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

func (s *ProjectService) ListAll(ctx context.Context, options *types.ListProjectsOptions) ([]*types.Project, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allProjects []*types.Project
	for {
		select {
		case <-ticker.C:
			projects, resp, err := s.List(ctx, options)
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

func (s *ProjectService) Get(ctx context.Context, id int64) (*types.Project, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var p *types.Project
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *ProjectService) GetByNameWithGroup(ctx context.Context, pathWithGroup string) (*types.Project, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%s", url.QueryEscape(pathWithGroup)), nil)
	if err != nil {
		return nil, nil, err
	}

	var p *types.Project
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

type Languages map[string]float64

func (s *ProjectService) GetLanguages(ctx context.Context, projectId int64) (Languages, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/languages", projectId), nil)
	if err != nil {
		return nil, nil, err
	}

	var p Languages
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}
