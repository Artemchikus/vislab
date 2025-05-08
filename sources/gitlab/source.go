package gitlab

import (
	"context"
	"vislab/config"
	"vislab/sources/gitlab/types"
)

type Source struct {
	gitClient *Client
	weight    int64
}

func NewSource(gitSourceConfig *config.GitSourceConfig) (*Source, error) {
	options := GetOptions(gitSourceConfig.Client)

	git, err := NewClient(gitSourceConfig.Client.Token, gitSourceConfig.Client.BaseURL, options...)
	if err != nil {
		return nil, err
	}

	s := &Source{
		gitClient: git,
		weight:    gitSourceConfig.Weight,
	}

	return s, nil
}

func (s *Source) GetData(ctx context.Context, serviceId int64) (*types.All, error) {
	project, _, err := s.gitClient.Projects.Get(ctx, serviceId)
	if err != nil {
		return nil, err
	}

	latestTag, _, err := s.gitClient.Tags.GetLatest(ctx, serviceId)
	if err != nil {
		return nil, err
	}

	all := &types.All{
		Project:   project,
		LatestTag: latestTag,
	}

	return all, nil
}

func (s *Source) Weight() int64 {
	return s.weight
}
