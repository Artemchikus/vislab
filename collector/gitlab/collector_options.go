package gitlabcollector

import (
	"fmt"
	"sort"
	"vislab/collector"
	gtlabjobsteps "vislab/collector/gitlab/steps"
	"vislab/sources/gitlab"
	"vislab/sources/migrations"
	"vislab/sources/yaml"
)

func WithYamlSource(yamlSource *yaml.Source, configPaths []string, fromGitlab bool) collector.CollectorOption {
	return func(c collector.Collector) error {
		collector, ok := c.(*Collector)
		if !ok {
			return fmt.Errorf("invalid collector type")
		}

		step := gtlabjobsteps.NewYamlStep(configPaths, collector.gitlabClient, yamlSource, fromGitlab)
		collector.steps = insertStepByWeight(collector.steps, step)
		return nil
	}
}

func WithGitlabSource(gitlabSource *gitlab.Source) collector.CollectorOption {
	return func(c collector.Collector) error {
		collector, ok := c.(*Collector)
		if !ok {
			return fmt.Errorf("invalid collector type")
		}

		step := gtlabjobsteps.NewGitlabStep(gitlabSource)
		collector.steps = insertStepByWeight(collector.steps, step)
		return nil
	}
}

func WithMigrationSource(migrationSource *migrations.Source, migrationsDirs []string) collector.CollectorOption {
	return func(c collector.Collector) error {
		collector, ok := c.(*Collector)
		if !ok {
			return fmt.Errorf("invalid collector type")
		}

		step := gtlabjobsteps.NewMigrationStep(migrationsDirs, collector.gitlabClient, migrationSource)
		collector.steps = insertStepByWeight(collector.steps, step)
		return nil
	}
}

func WithGitlabGroups(groups []string) collector.CollectorOption {
	return func(c collector.Collector) error {
		collector, ok := c.(*Collector)
		if !ok {
			return fmt.Errorf("invalid collector type")
		}

		collector.gitlabGroups = groups
		return nil
	}
}

func WithReleaseProject(project, releaseFile string, releaseTag string, releaseYamlSource *yaml.Source) collector.CollectorOption {
	return func(c collector.Collector) error {
		collector, ok := c.(*Collector)
		if !ok {
			return fmt.Errorf("invalid collector type")
		}

		collector.releaseProject = project
		collector.releaseFile = releaseFile
		collector.releaseTag = releaseTag
		collector.releaseYamlSource = releaseYamlSource
		return nil
	}
}

func WithReleaseTag(tag string) collector.CollectorOption {
	return func(c collector.Collector) error {
		collector, ok := c.(*Collector)
		if !ok {
			return fmt.Errorf("invalid collector type")
		}

		collector.releaseTag = tag
		return nil
	}
}

func insertStepByWeight(steps []gtlabjobsteps.Step, newStep gtlabjobsteps.Step) []gtlabjobsteps.Step {
	insertIdx := sort.Search(len(steps), func(i int) bool {
		return steps[i].Weight() > newStep.Weight()
	})

	steps = append(steps, nil)

	if insertIdx < len(steps)-1 {
		copy(steps[insertIdx+1:], steps[insertIdx:len(steps)-1])
	}

	steps[insertIdx] = newStep

	return steps
}
