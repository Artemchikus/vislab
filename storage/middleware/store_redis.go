package storefuncs

import (
	"context"
	"log/slog"
	"strings"
	"vislab/libs/check"
	"vislab/libs/ptr"
	"vislab/storage"
	storeTypes "vislab/storage/neo4j/types"
	"vislab/types"
)

func storeRedis(ctx context.Context, redis *types.Redis, serviceNode *storeTypes.ConnNode, storage storage.Storage) error {
	redisNode, err := storeRedisNode(ctx, redis, serviceNode, storage)
	if err != nil {
		return err
	}

	slog.Info("getting redis databases", "redis", redis.Host)
	existingDatabases, err := storage.Redis().GetDBs(ctx, redisNode.ID)
	if err != nil {
		return err
	}

	if len(redis.Databases) == 0 {
		slog.Info("creating dummy redis database", "redis", redis.Host)
		if err := storeDummyRedisDB(ctx, redisNode, storage, serviceNode, existingDatabases); err != nil {
			return err
		}

		return nil
	}

	for _, database := range redis.Databases {
		databaseNode, err := storeRedisDB(ctx, database, redisNode, existingDatabases, storage)
		if err != nil {
			slog.Error("failed to store redis database", "error", err)
			continue
		}

		slog.Info("getting redis namespaces", "redis", redis.Host, "database", database.Name)
		existingNamespaces, err := storage.Redis().GetNamespaces(ctx, databaseNode.ID)
		if err != nil {
			return err
		}

		if len(database.Namespaces) == 0 {
			slog.Info("creating dummy redis namespace", "redis", redis.Host, "database", database.Name)
			if err := storeDummyRedisNamespace(ctx, databaseNode, storage, serviceNode, existingNamespaces); err != nil {
				return err
			}

			continue
		}

		for _, namespace := range database.Namespaces {
			namespaceNode, err := storeRedisNamespace(ctx, namespace, databaseNode, existingNamespaces, storage)
			if err != nil {
				slog.Error("failed to store redis namespace", "error", err)
			}

			slog.Info("creating svc-namespace connection", "from_id", serviceNode.ID, "to_id", namespaceNode.ID, "type", storeTypes.ConnUses)
			if err := storage.Connection().Create(ctx, serviceNode, namespaceNode, storeTypes.ConnUses); err != nil {
				return err
			}
		}
	}

	return nil
}

func storeRedisNode(ctx context.Context, redis *types.Redis, serviceNode *storeTypes.ConnNode, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeRedis := &storeTypes.Redis{
		Host:   redis.Host,
		Port:   redis.Port,
		Master: redis.Master,
	}

	redisNode := &storeTypes.ConnNode{
		Class: storeTypes.RedisClass,
	}

	slog.Info("updating redis", "redis", redis.Host)
	dbRedis, err := storage.Redis().Update(ctx, storeRedis)
	if err == nil {
		redisNode.ID = *dbRedis.UID
		return redisNode, nil
	}
	slog.Error("error updating redis", "redis", redis.Host, "error", err)

	if strings.Contains(err.Error(), "nothing to update") {
		dbRedis, err := storage.Redis().Get(ctx, *storeRedis.Host)
		if err != nil {
			return nil, err
		}

		redisNode.ID = *dbRedis.UID
		return redisNode, nil
	}

	slog.Info("creating redis", "redis", redis.Host)
	id, err := storage.Redis().Create(ctx, storeRedis)
	if err != nil {
		return nil, err
	}

	redisNode.ID = id
	return redisNode, nil
}

func storeRedisDB(ctx context.Context, database *types.RedisDB, redisNode *storeTypes.ConnNode, existingDatabases []*storeTypes.RedisDB, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeDatabase := &storeTypes.RedisDB{
		Name: database.Name,
	}

	databaseNode := &storeTypes.ConnNode{
		Class: storeTypes.RedisDBClass,
	}

	for _, existingDatabase := range existingDatabases {
		if check.ComparePointers(existingDatabase.Name, database.Name) {
			if existingDatabase.Equal(storeDatabase) {
				databaseNode.ID = *existingDatabase.UID
				return databaseNode, nil
			}

			storeDatabase.UID = existingDatabase.UID

			slog.Info("updating rabbit db", "database", database.Name)
			dbDatabase, err := storage.Redis().UpdateDB(ctx, storeDatabase)
			if err != nil {
				return nil, err
			}

			databaseNode.ID = *dbDatabase.UID
			return databaseNode, nil
		}
	}

	slog.Info("creating redis database", "database", database.Name)
	id, err := storage.Redis().CreateDB(ctx, storeDatabase)
	if err != nil {
		return nil, err
	}

	databaseNode.ID = id

	slog.Info("creating redis-db connection", "from_id", databaseNode.ID, "to_id", redisNode.ID, "type", storeTypes.ConnIN)
	if err := storage.Connection().Create(ctx, databaseNode, redisNode, storeTypes.ConnIN); err != nil {
		return nil, err
	}

	return databaseNode, nil
}

func storeRedisNamespace(ctx context.Context, namespace *types.RedisNamespace, databaseNode *storeTypes.ConnNode, existingNamespaces []*storeTypes.RedisNamespace, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeNamespace := &storeTypes.RedisNamespace{
		Name: namespace.Name,
	}

	namespaceNode := &storeTypes.ConnNode{
		Class: storeTypes.RedisNSClass,
	}

	for _, existingNamespace := range existingNamespaces {
		if check.ComparePointers(existingNamespace.Name, namespace.Name) {
			if existingNamespace.Equal(storeNamespace) {
				namespaceNode.ID = *existingNamespace.UID
				return namespaceNode, nil
			}

			storeNamespace.UID = existingNamespace.UID

			slog.Info("updating redis namespace", "namespace", namespace.Name)
			dbNamespace, err := storage.Redis().UpdateNamespace(ctx, storeNamespace)
			if err != nil {
				return nil, err
			}

			namespaceNode.ID = *dbNamespace.UID
			return namespaceNode, nil
		}
	}

	slog.Info("creating redis namespace", "namespace", namespace.Name)
	id, err := storage.Redis().CreateNamespace(ctx, storeNamespace)
	if err != nil {
		return nil, err
	}

	namespaceNode.ID = id

	slog.Info("creating redisDB-namespace connection", "from_id", namespaceNode.ID, "to_id", databaseNode.ID, "type", storeTypes.ConnIN)

	if err := storage.Connection().Create(ctx, namespaceNode, databaseNode, storeTypes.ConnIN); err != nil {
		return nil, err
	}

	return namespaceNode, nil
}

func storeDummyRedisDB(ctx context.Context, redisNode *storeTypes.ConnNode, storage storage.Storage, serviceNode *storeTypes.ConnNode, existingDatabases []*storeTypes.RedisDB) error {
	dummyDB := &types.RedisDB{
		Name: ptr.Ptr("dummy"),
	}

	dbNode, err := storeRedisDB(ctx, dummyDB, redisNode, existingDatabases, storage)
	if err != nil {
		return err
	}

	if err := storeDummyRedisNamespace(ctx, dbNode, storage, serviceNode, []*storeTypes.RedisNamespace{}); err != nil {
		return err
	}

	return nil
}

func storeDummyRedisNamespace(ctx context.Context, namespaceNode *storeTypes.ConnNode, storage storage.Storage, serviceNode *storeTypes.ConnNode, existingNamespaces []*storeTypes.RedisNamespace) error {
	dummyNamespace := &types.RedisNamespace{
		Name: ptr.Ptr("dummy"),
	}

	namespaceNode, err := storeRedisNamespace(ctx, dummyNamespace, namespaceNode, existingNamespaces, storage)
	if err != nil {
		return err
	}

	slog.Info("creating svc-namespace connection", "from_id", serviceNode.ID, "to_id", namespaceNode.ID, "type", "dummy")
	if err := storage.Connection().Create(ctx, serviceNode, namespaceNode, "dummy"); err != nil {
		return err
	}

	return nil
}
