package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"vislab/libs/ptr"
	"vislab/sources/gitlab/types"
)

type FIleService struct {
	client *Client
}

func (s *FIleService) List(ctx context.Context, options *types.ListFilesOptions, projectId int64) ([]*types.ListFile, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/tree", projectId), options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.ListFile
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *FIleService) ListAll(ctx context.Context, options *types.ListFilesOptions, projectId int64) ([]*types.ListFile, *Response, error) {
	options.Page = 1
	options.Recursive = ptr.Ptr(true)

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allFiles []*types.ListFile
	for {
		select {
		case <-ticker.C:
			files, resp, err := s.List(ctx, options, projectId)
			if err != nil {
				return nil, resp, err
			}

			allFiles = append(allFiles, files...)

			if resp.NextPage == 0 {
				return allFiles, nil, nil
			}
			options.Page = int64(resp.NextPage)

		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}

func (s *FIleService) Get(ctx context.Context, path string, projectId int64, ref string) (*types.File, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/repository/files/%s", projectId, url.QueryEscape(path)), &types.ListFilesOptions{Ref: &ref})
	if err != nil {
		return nil, nil, err
	}

	var p *types.File
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *FIleService) ListDir(ctx context.Context, path string, projectId int64, ref string) ([]*types.ListFile, *Response, error) {
	options := &types.ListFilesOptions{Path: &path, Ref: &ref}

	return s.ListAll(ctx, options, projectId)
}
