package defaultaggregator

import (
	"context"
	yamlTypes "vislab/sources/yaml/types"
	"vislab/types"
)

func setPostgres(ctx context.Context, in []*yamlTypes.Postgresql, out *types.All) error {
	for _, postgresql := range in {
		newPostgres := &types.Postgresql{
			Host: postgresql.Host,
			Port: postgresql.Port,
			User: postgresql.User,
		}

		for _, database := range postgresql.Databases {
			newDB := &types.PostgresqlDB{
				Name: database.Name,
				// Owner: database.Owner,
				ForMigrations: database.ForMigrations,
			}

			for _, scheme := range database.Schemes {
				newScheme := &types.PostgresqlScheme{
					Name: scheme.Name,
					// Owner: scheme.Owner,
				}

				newDB.Schemes = append(newDB.Schemes, newScheme)
			}

			newPostgres.Databases = append(newPostgres.Databases, newDB)
		}

		out.Postgresqls = append(out.Postgresqls, newPostgres)
	}

	return nil
}
