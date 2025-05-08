package storefuncs

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"vislab/libs/check"
	"vislab/libs/ptr"
	"vislab/storage"
	storeTypes "vislab/storage/neo4j/types"
	"vislab/types"
)

func storePostgres(ctx context.Context, postgres *types.Postgresql, serviceNode *storeTypes.ConnNode, storage storage.Storage) error {
	postgresNode, err := storePostgresNode(ctx, postgres, storage)
	if err != nil {
		return err
	}

	slog.Info("getting postgres databases", "postgres", postgres.Host)
	existingDatabases, err := storage.Postgres().GetDBs(ctx, postgresNode.ID)
	if err != nil {
		return err
	}

	if len(postgres.Databases) == 0 {
		slog.Info("creating dummy postgres database", "postgres", postgres.Host)
		if err := storeDummyPostgresDB(ctx, postgresNode, storage, serviceNode, existingDatabases); err != nil {
			return err
		}

		return nil
	}

	for _, database := range postgres.Databases {
		databaseNode, err := storePostgresDB(ctx, database, postgresNode, existingDatabases, storage)
		if err != nil {
			slog.Error("failed to store postgres database", "error", err)
		}

		slog.Info("getting postgres schemes", "database", database.Name, "postgres", postgres.Host)
		existingSchemes, err := storage.Postgres().GetSchemes(ctx, databaseNode.ID)
		if err != nil {
			return err
		}

		if len(database.Schemes) == 0 {
			slog.Info("creating dummy postgres scheme", "database", database.Name, "postgres", postgres.Host)
			if err := storePublicScheme(ctx, databaseNode, storage, serviceNode, existingSchemes); err != nil {
				return err
			}

			continue
		}

		for _, scheme := range database.Schemes {
			schemeNode, err := storePostgresScheme(ctx, scheme, databaseNode, existingSchemes, storage)
			if err != nil {
				slog.Error("failed to store postgres scheme", "error", err)
				continue
			}

			slog.Info("getting postgres tables", "scheme", scheme.Name, "database", database.Name, "postgres", postgres.Host)
			existingTables, err := storage.Postgres().GetTables(ctx, schemeNode.ID)
			if err != nil {
				return err
			}

			if len(scheme.Tables) == 0 {
				slog.Info("creating dummy postgres table", "scheme", scheme.Name, "database", database.Name, "postgres", postgres.Host)
				if err := storeDummyPostgresTable(ctx, schemeNode, storage, serviceNode, existingTables); err != nil {
					return err
				}

				continue
			}

			for _, table := range scheme.Tables {
				tableNode, err := storePostgresTable(ctx, table, schemeNode, existingTables, storage)
				if err != nil {
					slog.Error("failed to store postgres table", "error", err)
					continue
				}

				slog.Info("creating svc-table connection", "from_id", serviceNode.ID, "to_id", tableNode.ID, "type", storeTypes.ConnUses)

				if err := storage.Connection().Create(ctx, serviceNode, tableNode, storeTypes.ConnUses); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func storePostgresNode(ctx context.Context, postgres *types.Postgresql, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storePostgresql := &storeTypes.Postgresql{
		Host: postgres.Host,
		Port: postgres.Port,
		User: postgres.User,
	}

	postgresNode := &storeTypes.ConnNode{
		Class: storeTypes.PostgresClass,
	}

	slog.Info("updating postgres", "postgres", postgres.Host)
	dbPostgres, err := storage.Postgres().Update(ctx, storePostgresql)
	if err == nil {
		postgresNode.ID = *dbPostgres.UID
		return postgresNode, nil
	}
	slog.Error("error updating postgres", "postgres", postgres.Host, "error", err)

	if strings.Contains(err.Error(), "nothing to update") {
		dbPostgres, err := storage.Postgres().Get(ctx, *storePostgresql.Host)
		if err != nil {
			return nil, err
		}

		postgresNode.ID = *dbPostgres.UID
		return postgresNode, nil
	}

	slog.Info("creating postgres", "postgres", postgres.Host)
	id, err := storage.Postgres().Create(ctx, storePostgresql)
	if err != nil {
		return nil, err
	}

	postgresNode.ID = id
	return postgresNode, nil
}

func storePostgresDB(ctx context.Context, db *types.PostgresqlDB, postgresNode *storeTypes.ConnNode, existingDBs []*storeTypes.PostgresqlDB, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeDB := &storeTypes.PostgresqlDB{
		Name: db.Name,
	}

	dbNode := &storeTypes.ConnNode{
		Class: storeTypes.PostgresDBClass,
	}

	for _, existingDB := range existingDBs {
		if check.ComparePointers(db.Name, existingDB.Name) {
			if existingDB.Equal(storeDB) {
				dbNode.ID = *existingDB.UID
				return dbNode, nil
			}

			storeDB.UID = existingDB.UID

			slog.Info("updating postgres database", "database", db.Name)
			dbDatabase, err := storage.Postgres().UpdateDB(ctx, storeDB)
			if err != nil {
				return nil, err
			}

			dbNode.ID = *dbDatabase.UID
			return dbNode, nil
		}
	}

	slog.Info("creating postgres database", "database", db.Name)
	id, err := storage.Postgres().CreateDB(ctx, storeDB)
	if err != nil {
		return nil, err
	}

	dbNode.ID = id

	slog.Info("creating db-postgres connection", "from_id", dbNode.ID, "to_id", postgresNode.ID, "type", storeTypes.ConnIN)
	if err := storage.Connection().Create(ctx, dbNode, postgresNode, storeTypes.ConnIN); err != nil {
		return nil, fmt.Errorf("failed to create db-postgres connection: %w", err)
	}

	return dbNode, nil
}

func storePostgresScheme(ctx context.Context, scheme *types.PostgresqlScheme, dbNode *storeTypes.ConnNode, existingSchemes []*storeTypes.PostgresqlScheme, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeScheme := &storeTypes.PostgresqlScheme{
		Name: scheme.Name,
	}

	schemeNode := &storeTypes.ConnNode{
		Class: storeTypes.PostgresSchemeClass,
	}

	for _, existingScheme := range existingSchemes {
		if check.ComparePointers(scheme.Name, existingScheme.Name) {
			if existingScheme.Equal(storeScheme) {
				schemeNode.ID = *existingScheme.UID
				return schemeNode, nil
			}

			storeScheme.UID = existingScheme.UID

			slog.Info("updating postgres scheme", "scheme", scheme.Name)
			dbScheme, err := storage.Postgres().UpdateScheme(ctx, storeScheme)
			if err != nil {
				return nil, err
			}

			schemeNode.ID = *dbScheme.UID
			return schemeNode, nil
		}
	}

	slog.Info("creating postgres scheme", "scheme", scheme.Name)
	id, err := storage.Postgres().CreateScheme(ctx, storeScheme)
	if err != nil {
		return nil, err
	}

	schemeNode.ID = id

	slog.Info("creating scheme-db connection", "from_id", schemeNode.ID, "to_id", dbNode.ID, "type", storeTypes.ConnIN)
	if err := storage.Connection().Create(ctx, schemeNode, dbNode, storeTypes.ConnIN); err != nil {
		return nil, fmt.Errorf("failed to create scheme-db connection: %w", err)
	}

	return schemeNode, nil
}

func storePostgresTable(ctx context.Context, table *types.PostgresqlTable, schemeNode *storeTypes.ConnNode, existingTables []*storeTypes.PostgresqlTable, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeTable := &storeTypes.PostgresqlTable{
		Name: table.Name,
	}

	tableNode := &storeTypes.ConnNode{
		Class: storeTypes.PostgresTableClass,
	}

	for _, existingTable := range existingTables {
		if check.ComparePointers(table.Name, existingTable.Name) {
			if existingTable.Equal(storeTable) {
				tableNode.ID = *existingTable.UID
				return tableNode, nil
			}

			storeTable.UID = existingTable.UID

			slog.Info("updating postgres table", "table", table.Name)
			dbTable, err := storage.Postgres().UpdateTable(ctx, storeTable)
			if err != nil {
				return nil, err
			}

			tableNode.ID = *dbTable.UID
			return tableNode, nil
		}
	}

	slog.Info("creating postgres table", "table", table.Name)
	id, err := storage.Postgres().CreateTable(ctx, storeTable)
	if err != nil {
		return nil, err
	}

	tableNode.ID = id

	slog.Info("creating table-scheme connection", "from_id", tableNode.ID, "to_id", schemeNode.ID, "type", storeTypes.ConnIN)
	if err := storage.Connection().Create(ctx, tableNode, schemeNode, storeTypes.ConnIN); err != nil {
		return nil, fmt.Errorf("failed to create table-scheme connection: %w", err)
	}

	return tableNode, nil
}

func storeDummyPostgresDB(ctx context.Context, postgresNode *storeTypes.ConnNode, storage storage.Storage, serviceNode *storeTypes.ConnNode, existingDatabases []*storeTypes.PostgresqlDB) error {
	dummyDB := &types.PostgresqlDB{
		Name: ptr.Ptr("dummy"),
	}

	dbNode, err := storePostgresDB(ctx, dummyDB, postgresNode, existingDatabases, storage)
	if err != nil {
		return err
	}

	if err := storePublicScheme(ctx, dbNode, storage, serviceNode, []*storeTypes.PostgresqlScheme{}); err != nil {
		return err
	}

	return nil
}

func storePublicScheme(ctx context.Context, postgresDB *storeTypes.ConnNode, storage storage.Storage, serviceNode *storeTypes.ConnNode, existingSchemes []*storeTypes.PostgresqlScheme) error {
	publicScheme := &types.PostgresqlScheme{
		Name: ptr.Ptr("public"),
	}

	schemeNode, err := storePostgresScheme(ctx, publicScheme, postgresDB, existingSchemes, storage)
	if err != nil {
		return err
	}

	if err := storeDummyPostgresTable(ctx, schemeNode, storage, serviceNode, []*storeTypes.PostgresqlTable{}); err != nil {
		return err
	}

	return nil
}

func storeDummyPostgresTable(ctx context.Context, schemeNode *storeTypes.ConnNode, storage storage.Storage, serviceNode *storeTypes.ConnNode, existingTables []*storeTypes.PostgresqlTable) error {
	dummyTable := &types.PostgresqlTable{
		Name: ptr.Ptr("dummy"),
	}

	tableNode, err := storePostgresTable(ctx, dummyTable, schemeNode, existingTables, storage)
	if err != nil {
		return err
	}

	slog.Info("creating svc-table connection", "from_id", serviceNode.ID, "to_id", tableNode.ID, "type", "dummy")
	if err := storage.Connection().Create(ctx, serviceNode, tableNode, "dummy"); err != nil {
		return err
	}

	return nil
}
