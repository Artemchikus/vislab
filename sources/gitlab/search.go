package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vislab/libs/ptr"
	"vislab/sources/gitlab/types"
)

type SearchService struct {
	client *Client
}

func (s *SearchService) Search(ctx context.Context, options *types.SearchOptions, projectId int64) ([]*types.SearchResult, *Response, error) {
	if options.PerPage == 0 {
		options.PerPage = s.client.perPage
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("projects/%d/search", projectId), options)
	if err != nil {
		return nil, nil, err
	}

	var p []*types.SearchResult
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, nil
}

func (s *SearchService) SearchFiles(ctx context.Context, query string, projectId int64, options *types.SearchOptions) ([]*types.SearchResult, *Response, error) {
	options.Search = ptr.Ptr(query)
	options.Scope = ptr.Ptr("blobs")
	return s.Search(ctx, options, projectId)
}

func (s *SearchService) SearchAllFiles(ctx context.Context, query string, projectId int64, options *types.SearchOptions) ([]*types.SearchResult, *Response, error) {
	options.Page = 1

	ticker := time.NewTicker(s.client.rateLimit)
	defer ticker.Stop()

	var allFiles []*types.SearchResult
	for {
		select {
		case <-ticker.C:
			files, resp, err := s.SearchFiles(ctx, query, projectId, options)
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

func (s *SearchService) GenerateFileQuery(line string, options *types.FileQuery) string {
	query := line

	if options.Path != "" {
		query = fmt.Sprintf("%s path:%s", query, options.Path)
	}
	if options.Filename != "" {
		query = fmt.Sprintf("%s filename:%s", query, options.Filename)
	}
	if options.Extension != "" {
		query = fmt.Sprintf("%s extension:%s", query, options.Extension)
	}

	return query
}
