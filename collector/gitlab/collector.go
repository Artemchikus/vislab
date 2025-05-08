package gitlabcollector

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	defaultaggregator "vislab/aggregator/default"
	"vislab/collector"
	gtlabjobsteps "vislab/collector/gitlab/steps"
	"vislab/config"
	"vislab/sources/gitlab"
	"vislab/sources/gitlab/types"
	"vislab/sources/migrations"
	"vislab/sources/yaml"
	yamlTypes "vislab/sources/yaml/types"
	"vislab/storage"
	storefuncs "vislab/storage/middleware"
)

type Collector struct {
	gitlabGroups   []string
	releaseProject string
	releaseFile    string
	releaseTag     string

	steps []gtlabjobsteps.Step

	releaseYamlSource *yaml.Source
	gitlabClient      *gitlab.Client
	storage           storage.Storage
}

func New(gitlabClient *gitlab.Client, storage storage.Storage, options ...collector.CollectorOption) (*Collector, error) {
	if gitlabClient == nil {
		return nil, fmt.Errorf("no gitlab client specified")
	}
	if storage == nil {
		return nil, fmt.Errorf("no storage specified")
	}

	gitlabCollector := &Collector{
		gitlabClient: gitlabClient,
		storage:      storage,
		steps:        []gtlabjobsteps.Step{},
	}

	for _, option := range options {
		if err := option(gitlabCollector); err != nil {
			return nil, err
		}
	}

	if len(gitlabCollector.steps) == 0 {
		return nil, fmt.Errorf("no source specified")
	}

	return gitlabCollector, nil
}

func (c *Collector) Collect(ctx context.Context) error {
	switch {
	case c.releaseProject != "" && c.releaseFile != "":
		return c.collectFromReleaseFile(ctx)
	default:
		return c.collectAll(ctx)
	}
}

func (c *Collector) Update(ctx context.Context, options ...collector.CollectorOption) error {
	for _, option := range options {
		if err := option(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Collector) collectAll(ctx context.Context) error {
	neededGroups, err := getNeededGroups(c.gitlabGroups, c.gitlabClient)
	if err != nil {
		return fmt.Errorf("failed to get needed groups: %w", err)
	}

	neededProjects, err := getNeededProjects(neededGroups, c.gitlabClient)
	if err != nil {
		return fmt.Errorf("failed to get needed projects: %w", err)
	}

	for _, project := range neededProjects {
		params := &gtlabjobsteps.StepParams{
			ServiceId:  *project.ID,
			ServiceRef: *project.DefaultBranch,
		}

		if err := c.collectProject(ctx, params); err != nil {
			slog.Error("failed to collect project data", "err", err, "project", project.Name)
			continue
		}
	}
	return nil
}

func (c *Collector) collectFromReleaseFile(ctx context.Context) error {
	project, _, err := c.gitlabClient.Projects.GetByNameWithGroup(ctx, c.releaseProject)
	if err != nil {
		return fmt.Errorf("failed to get release project: %w", err)
	}

	releaseFile64, _, err := c.gitlabClient.Files.Get(ctx, c.releaseFile, *project.ID, c.releaseTag)
	if err != nil {
		return fmt.Errorf("failed to get release file: %w", err)
	}

	releaseFile, err := base64.StdEncoding.DecodeString(releaseFile64.Content)
	if err != nil {
		return fmt.Errorf("failed to decode release file: %w", err)
	}

	releaseInfo := &yamlTypes.All{}

	if err := c.releaseYamlSource.GetData(ctx, releaseFile, releaseInfo); err != nil {
		return fmt.Errorf("failed to get data from release file: %w", err)
	}

	for _, service := range releaseInfo.Service.Instances {
		var project *types.Project
		if service.ProjectID == nil {
			project, _, err = c.gitlabClient.Projects.GetByNameWithGroup(ctx, *service.FullName)
			if err != nil {
				slog.Error("failed to get project", "err", err, "project", *service.FullName)
				continue
			}
		} else {
			project, _, err = c.gitlabClient.Projects.Get(ctx, *service.ProjectID)
			if err != nil {
				slog.Error("failed to get project", "err", err, "project", *service.ProjectID)
				continue
			}
		}

		params := &gtlabjobsteps.StepParams{
			ServiceId:  *project.ID,
			ServiceRef: *service.Tag,
		}

		if err := c.collectProject(ctx, params); err != nil {
			slog.Error("failed to collect project data", "err", err, "project", project.Name)
			continue
		}
	}

	return nil
}

func (c *Collector) collectProject(ctx context.Context, params *gtlabjobsteps.StepParams) error {
	aggr, err := defaultaggregator.New()
	if err != nil {
		return fmt.Errorf("failed to create aggregator: %w", err)
	}

	params.Aggregator = aggr

	for _, step := range c.steps {
		if err := step.Run(ctx, params); err != nil {
			return fmt.Errorf("failed to run step: %w", err)
		}
	}

	aggrData, err := aggr.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to aggregate data: %w", err)
	}

	exist, err := isAlreadyExist(ctx, aggrData.Service, c.storage)
	if err != nil {
		return fmt.Errorf("failed to check if already exist: %w", err)
	}

	if exist {
		slog.Debug("project already exist", "service", aggrData.Service)
		return nil
	}

	if err := storefuncs.StoreResources(ctx, aggrData, c.storage); err != nil {
		return fmt.Errorf("failed to store resource: %w", err)
	}

	return nil
}

func GetOptions(collectorConf *config.CollectorConfig, sourcesConf *config.SourcesConfig) ([]collector.CollectorOption, error) {
	options := []collector.CollectorOption{}

	if collectorConf.GitLab.ReleaseProject != nil {
		slog.Info("release project enabled")
		releaseYamlSource, err := yaml.NewSource(&config.YamlSourceConfig{
			ParseConfigPath: collectorConf.GitLab.ReleaseProject.ParseConfigPath,
			Weight:          0,
			FromGitlab:      true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create release yaml source: %w", err)
		}
		options = append(options, WithReleaseProject(collectorConf.GitLab.ReleaseProject.Project, collectorConf.GitLab.ReleaseProject.ReleaseFilePath, collectorConf.GitLab.ReleaseProject.Tag, releaseYamlSource))
	}
	if collectorConf.GitLab.Groups != nil {
		slog.Info("groups filter enabled")
		options = append(options, WithGitlabGroups(collectorConf.GitLab.Groups))
	}
	if sourcesConf.Migration != nil {
		slog.Info("migration source enabled")
		migrationSource, err := migrations.NewSource(sourcesConf.Migration)
		if err != nil {
			return nil, fmt.Errorf("failed to create migration source: %w", err)
		}
		options = append(options, WithMigrationSource(migrationSource, collectorConf.MigrationPaths))
	}
	if sourcesConf.Yaml != nil {
		slog.Info("yaml source enabled")
		yamlSource, err := yaml.NewSource(sourcesConf.Yaml)
		if err != nil {
			return nil, fmt.Errorf("failed to create yaml source: %w", err)
		}
		options = append(options, WithYamlSource(yamlSource, collectorConf.ServiceConfigPaths, false))
	}
	if sourcesConf.GitLab != nil {
		slog.Info("gitlab source enabled")
		gitlabSource, err := gitlab.NewSource(sourcesConf.GitLab)
		if err != nil {
			return nil, fmt.Errorf("failed to create gitlab source: %w", err)
		}
		options = append(options, WithGitlabSource(gitlabSource))
	}

	return options, nil
}
