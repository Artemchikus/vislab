package gtlabjobsteps

import (
	"context"
	"fmt"
	"log/slog"
	"vislab/sources/gitlab"
)

type GitlabStep struct {
	gitlabSource *gitlab.Source
}

func NewGitlabStep(gitlabSource *gitlab.Source) *GitlabStep {
	return &GitlabStep{
		gitlabSource: gitlabSource,
	}
}

func (s *GitlabStep) Run(ctx context.Context, params *StepParams) error {
	slog.Info("running gitlab step", "service_id", params.ServiceId, "ref", params.ServiceRef)
	gitlabSourceData, err := s.gitlabSource.GetData(ctx, params.ServiceId)
	if err != nil {
		return fmt.Errorf("failed to get data from gitlab source: %w", err)
	}

	gitlabSourceData.LatestTag.Name = &params.ServiceRef

	if err := params.Aggregator.Set(ctx, gitlabSourceData); err != nil {
		return fmt.Errorf("failed to set gitlab source data: %w", err)
	}

	return nil
}

func (s *GitlabStep) Weight() int64 {
	return s.gitlabSource.Weight()
}
