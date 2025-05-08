package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Storage   *StorageConfig   `yaml:"storage"`
		Sources   *SourcesConfig   `yaml:"sources"`
		Collector *CollectorConfig `yaml:"collector"`
		Updater   *UpdaterConfig   `yaml:"updater"`
	}
	UpdaterConfig struct {
		Port string `yaml:"port"`
	}
	CollectorConfig struct {
		ParallelJobs       int64                  `yaml:"parallel_jobs"`
		ServiceConfigPaths []string               `yaml:"service_config_paths"`
		MigrationPaths     []string               `yaml:"migration_paths"`
		GitLab             *GitLabCollectorConfig `yaml:"gitlab"`
	}
	GitLabCollectorConfig struct {
		Client         *GitLabClientConfig   `yaml:"client"`
		Groups         []string              `yaml:"groups"`
		ReleaseProject *ReleaseProjectConfig `yaml:"release_project"`
	}
	ReleaseProjectConfig struct {
		Project         string `yaml:"project"`
		ReleaseFilePath string `yaml:"release_file_path"`
		Tag             string `yaml:"tag"`
		ParseConfigPath string `yaml:"parse_config_path"`
	}
	StorageConfig struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
	}
	SourcesConfig struct {
		Yaml      *YamlSourceConfig      `yaml:"yaml"`
		GitLab    *GitSourceConfig       `yaml:"gitlab"`
		Migration *MigrationSourceConfig `yaml:"migration"`
	}
	YamlSourceConfig struct {
		ParseConfigPath string `yaml:"parse_config_path"`
		FromGitlab      bool   `yaml:"from_gitlab"`
		Weight          int64  `yaml:"weight"`
	}
	GitLabClientConfig struct {
		Token           string `yaml:"token"`
		BaseURL         string `yaml:"base_url"`
		UseArchived     bool   `yaml:"use_archived"`
		GitlabPerPage   int64  `yaml:"per_page"`
		GitlabTimeout   int64  `yaml:"timeout"`
		GitLabRateLimit int64  `yaml:"rate_limit"`
		GitlabAPIPrefix string `yaml:"api_prefix"`
	}
	MigrationSourceConfig struct {
		Weight int64 `yaml:"weight"`
	}
	GitSourceConfig struct {
		Client *GitLabClientConfig `yaml:"client"`
		Weight int64               `yaml:"weight"`
	}
)

func Get(confFile string) (config *Config, err error) {
	rawData, err := os.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	config = &Config{}

	if err := config.Unmarshal(rawData); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Unmarshal(data []byte) (err error) {
	c.Storage = &StorageConfig{}

	if err = yaml.Unmarshal(data, c); err != nil {
		return err
	}

	if c.Storage.Host == "" {
		return fmt.Errorf("config: Storage.host is empty")
	}
	if c.Storage.Port == "" {
		return fmt.Errorf("config: Storage.port is empty")
	}
	if c.Storage.User == "" {
		return fmt.Errorf("config: Storage.user is empty")
	}
	if c.Storage.Name == "" {
		return fmt.Errorf("config: Storage.name is empty")
	}
	if c.Storage.Password == "" {
		return fmt.Errorf("config: Storage.password is empty")
	}

	return nil
}
