package defaultaggregator

import (
	"context"
	"fmt"
	"log/slog"
	gitlabTypes "vislab/sources/gitlab/types"
	migrationTypes "vislab/sources/migrations/types"
	yamlTypes "vislab/sources/yaml/types"
	"vislab/types"
)

type Aggregator struct {
	data *types.All
}

func New() (*Aggregator, error) {
	myAggr := &Aggregator{
		data: &types.All{
			Service:       &types.Service{},
			Kafkas:        []*types.Kafka{},
			Redises:       []*types.Redis{},
			Postgresqls:   []*types.Postgresql{},
			RabbitMQs:     []*types.RabbitMQ{},
			OtherServices: []*types.Service{},
		},
	}

	return myAggr, nil
}

func (a *Aggregator) Get(ctx context.Context) (*types.All, error) {
	return a.data, nil
}

func (a *Aggregator) Set(ctx context.Context, data any) error {
	switch d := data.(type) {
	case *gitlabTypes.All:
		if err := a.setGitlab(ctx, d); err != nil {
			return err
		}
	case *yamlTypes.All:
		if err := a.setYaml(ctx, d); err != nil {
			return err
		}
	case *migrationTypes.All:
		if err := a.setMigration(ctx, d); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown data type: %T", data)
	}
	return nil
}

func (a *Aggregator) setGitlab(ctx context.Context, data *gitlabTypes.All) error {
	a.data.Service.Description = data.Project.Description
	a.data.Service.Name = data.Project.Name
	a.data.Service.Group = data.Project.Group.Name
	a.data.Service.LatestTag = data.LatestTag.Name
	a.data.Service.FullName = data.Project.PathWithGroup
	// a.data.Service.Language = data.Project.Language
	a.data.Service.Link = data.Project.WebURL
	a.data.Service.MainBranch = data.Project.DefaultBranch

	return nil
}

func (a *Aggregator) setYaml(ctx context.Context, data *yamlTypes.All) error {
	if data.Service != nil {
		if err := setService(ctx, data.Service.Instances, a.data); err != nil {
			return fmt.Errorf("failed to set service data: %w", err)
		}
	}

	if data.Kafka != nil {
		if err := setKafka(ctx, data.Kafka.Instances, a.data); err != nil {
			return fmt.Errorf("failed to set kafka data: %w", err)
		}
	}

	if data.Redis != nil {
		if err := setRedis(ctx, data.Redis.Instances, a.data); err != nil {
			return fmt.Errorf("failed to set redis data: %w", err)
		}
	}

	if data.Postgresql != nil {
		if err := setPostgres(ctx, data.Postgresql.Instances, a.data); err != nil {
			return fmt.Errorf("failed to set postgresql data: %w", err)
		}
	}

	if data.RabbitMQ != nil {
		if err := setRabbit(ctx, data.RabbitMQ.Instances, a.data); err != nil {
			return fmt.Errorf("failed to set rabbitmq data: %w", err)
		}
	}

	if data.OtherService != nil {
		if err := setOtherService(ctx, data.OtherService.Instances, a.data); err != nil {
			return fmt.Errorf("failed to set other services data: %w", err)
		}
	}
	return nil
}

func (a *Aggregator) setMigration(ctx context.Context, data *migrationTypes.All) error {
	if len(a.data.Postgresqls) == 0 {
		return fmt.Errorf("no postgresqls for migrations found")
	}
	migrationPostgres := a.data.Postgresqls[0]

	if len(migrationPostgres.Databases) == 0 {
		return fmt.Errorf("no databases for migrations found")
	}
	migrationDatabase := migrationPostgres.Databases[0]

Postgreses:
	for _, postgres := range a.data.Postgresqls {
		for _, database := range postgres.Databases {
			if database.ForMigrations == nil {
				continue
			}
			if *database.ForMigrations {
				migrationDatabase = database
				break Postgreses
			}
		}
	}

Tables:
	for _, table := range data.Tables {
		newTable := &types.PostgresqlTable{
			Name: &table.Name,
			// Owner: &table.Owner,
			Type: &table.Type,
		}

		if table.Type == "partition" { // TODO: make optional
			continue Tables
		}

		for _, scheme := range migrationDatabase.Schemes {
			if *scheme.Name == table.Schema {
				scheme.Tables = append(scheme.Tables, newTable)
				continue Tables
			}
		}
		slog.Warn("table schema not found", "table", table.Name, "schema", table.Schema)
		migrationDatabase.Schemes = append(migrationDatabase.Schemes, &types.PostgresqlScheme{
			Name: &table.Schema,
			Tables: []*types.PostgresqlTable{
				newTable,
			},
		})
	}

	return nil
}
