package gtlabjobsteps

import (
	"context"
	"encoding/base64"
	"log/slog"
	"vislab/sources/gitlab"
	"vislab/sources/migrations"
	migrationsTypes "vislab/sources/migrations/types"
)

type MigrationStep struct {
	migrationDirs   []string
	gitlabClient    *gitlab.Client
	migrationSource *migrations.Source
}

func NewMigrationStep(migrationDirs []string, gitlabClient *gitlab.Client, migrationSource *migrations.Source) *MigrationStep {
	return &MigrationStep{
		migrationDirs:   migrationDirs,
		gitlabClient:    gitlabClient,
		migrationSource: migrationSource,
	}
}

func (s *MigrationStep) Run(ctx context.Context, params *StepParams) error {
	for _, migrationDir := range s.migrationDirs {
		slog.Info("getting migration files", "service_id", params.ServiceId, "ref", params.ServiceRef, "path", migrationDir)
		migrationFiles, _, err := s.gitlabClient.Files.ListDir(ctx, migrationDir, params.ServiceId, params.ServiceRef)
		if err != nil {
			slog.Error("failed to get migration files", "err", err, "path", migrationDir, "service_id", params.ServiceId, "ref", params.ServiceRef)
			continue
		}

		all := &migrationsTypes.All{
			Tables:   make(map[string]*migrationsTypes.Table),
			Funcs:    make(map[string][]*migrationsTypes.Func),
			Indexes:  make(map[string]*migrationsTypes.Index),
			Triggers: make(map[string]*migrationsTypes.Trigger),
			Types:    make(map[string]*migrationsTypes.Type),
		}

		for _, migrationFile := range migrationFiles {
			migrationData64, _, err := s.gitlabClient.Files.Get(ctx, migrationFile.Path, params.ServiceId, params.ServiceRef)
			if err != nil {
				slog.Error("failed to get migration file", "err", err, "path", migrationFile.Path, "service_id", params.ServiceId, "ref", params.ServiceRef)
				continue
			}

			migrationData, err := base64.StdEncoding.DecodeString(migrationData64.Content)
			if err != nil {
				slog.Error("failed to decode migration file", "err", err, "path", migrationFile.Path, "service_id", params.ServiceId, "ref", params.ServiceRef)
				continue
			}

			if err := s.migrationSource.GetData(ctx, migrationData, all); err != nil {
				slog.Error("failed to get data from migration file", "err", err, "path", migrationFile.Path, "service_id", params.ServiceId, "ref", params.ServiceRef)
				continue
			}
		}

		if err := params.Aggregator.Set(ctx, all); err != nil {
			slog.Error("failed to set migration files", "err", err, "path", migrationDir, "service_id", params.ServiceId, "ref", params.ServiceRef)
			continue
		}

		break // TODO: add support for multiple migration dirs (maybe)
	}
	return nil
}

func (s *MigrationStep) Weight() int64 {
	return s.migrationSource.Weight()
}
