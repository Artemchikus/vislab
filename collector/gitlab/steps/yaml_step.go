package gtlabjobsteps

import (
	"context"
	"encoding/base64"
	"log/slog"
	"os"
	"vislab/sources/gitlab"
	"vislab/sources/yaml"
	yamlTypes "vislab/sources/yaml/types"
)

type YamlStep struct {
	filePaths    []string
	gitlabClient *gitlab.Client
	yamlSource   *yaml.Source
	fromGitlab   bool
}

func NewYamlStep(filePaths []string, gitlabClient *gitlab.Client, yamlSource *yaml.Source, fromGitlab bool) *YamlStep {
	return &YamlStep{
		filePaths:    filePaths,
		gitlabClient: gitlabClient,
		yamlSource:   yamlSource,
		fromGitlab:   fromGitlab,
	}
}

func (s *YamlStep) Run(ctx context.Context, params *StepParams) error {
	for _, configPath := range s.filePaths {
		slog.Info("running yaml step", "service_id", params.ServiceId, "ref", params.ServiceRef, "path", configPath)

		var configData []byte

		if s.fromGitlab {
			config, _, err := s.gitlabClient.Files.Get(ctx, configPath, params.ServiceId, params.ServiceRef)
			if err != nil {
				slog.Error("failed to get config file", "err", err, "path", configPath, "service_id", params.ServiceId, "ref", params.ServiceRef)
				continue
			}

			configData, err = base64.StdEncoding.DecodeString(config.Content)
			if err != nil {
				slog.Error("failed to decode config file", "err", err, "path", configPath, "service_id", params.ServiceId, "ref", params.ServiceRef)
				continue
			}
		} else {
			service, _, err := s.gitlabClient.Projects.Get(ctx, params.ServiceId)
			if err != nil {
				slog.Error("failed to get service", "err", err, "service_id", params.ServiceId)
				continue
			}

			filePath := *service.PathWithGroup

			if len(s.filePaths) != 0 {
				filePath = s.filePaths[0] + "/" + *service.PathWithGroup + ".yaml"
			}

			configData, err = os.ReadFile(filePath)
			if err != nil {
				slog.Error("failed to read service config", "err", err, "service_id", params.ServiceId, "path", *service.PathWithGroup)
				continue
			}
		}

		all := &yamlTypes.All{}

		if err := s.yamlSource.GetData(ctx, configData, all); err != nil {
			slog.Error("failed to get data from config file", "err", err, "path", configPath, "service_id", params.ServiceId, "ref", params.ServiceRef)
			continue
		}

		if err := params.Aggregator.Set(ctx, all); err != nil {
			slog.Error("failed to set config file", "err", err, "path", configPath, "service_id", params.ServiceId, "ref", params.ServiceRef)
			continue
		}

		break // TODO: add support for multiple config files (maybe)
	}

	return nil
}

func (s *YamlStep) Weight() int64 {
	return s.yamlSource.Weight()
}
