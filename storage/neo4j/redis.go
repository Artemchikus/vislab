package neo4j

import (
	"context"
	"fmt"
	"strings"
	"vislab/storage"
	"vislab/storage/neo4j/types"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type (
	neo4jRedisRepo struct {
		db neo4j.DriverWithContext
	}
)

func (n *Neo4jStorage) Redis() storage.RedisRepository {
	if n.redisRepo != nil {
		return n.redisRepo
	}

	n.redisRepo = &neo4jRedisRepo{db: n.db}
	return n.redisRepo
}

func (n *neo4jRedisRepo) Create(ctx context.Context, redis *types.Redis) (string, error) {
	query := `CREATE
	(r:Redis {
		host: $host,
		port: $port,
		master: $master
	})
	RETURN r
	`

	args := map[string]any{
		"host":   redis.Host,
		"port":   redis.Port,
		"master": redis.Master,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("redis not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "r")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jRedisRepo) CreateDB(ctx context.Context, db *types.RedisDB) (string, error) {
	query := `CREATE
	(rd:RedisDB {
		name: $name
	})
	RETURN rd
	`

	args := map[string]any{
		"name": db.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("redis db node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "rd")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jRedisRepo) CreateNamespace(ctx context.Context, redisNS *types.RedisNamespace) (string, error) {
	query := `CREATE
	(rn:RedisNS {
		name: $name
	})
	RETURN rn
	`

	args := map[string]any{
		"name": redisNS.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("redis ns node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "rn")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jRedisRepo) Get(ctx context.Context, host string) (*types.Redis, error) {
	query := `MATCH
	(r:Redis)
	WHERE r.host = $host
	RETURN r
	`

	args := map[string]any{
		"host": host,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("redis not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "r")
	if err != nil {
		return nil, err
	}

	redis := &types.Redis{
		UID: &itemNode.ElementId,
	}

	if hostAny, ok := itemNode.Props["host"]; ok {
		host := hostAny.(string)
		redis.Host = &host
	}
	if portAny, ok := itemNode.Props["port"]; ok {
		port := portAny.(int64)
		redis.Port = &port
	}

	return redis, nil
}

func (n *neo4jRedisRepo) GetDBs(ctx context.Context, redisUID string) ([]*types.RedisDB, error) {
	query := `MATCH
	(rd:RedisDB)-[:IN]->(r:Redis)
	WHERE elementId(r) = $uid
	RETURN rd
	`

	args := map[string]any{
		"uid": redisUID,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	var dbs []*types.RedisDB

	if len(res.Records) == 0 {
		return dbs, nil
	}

	for _, record := range res.Records {
		itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](record, "rd")
		if err != nil {
			return nil, err
		}

		db := &types.RedisDB{
			UID: &itemNode.ElementId,
		}

		if nameAny, ok := itemNode.Props["name"]; ok {
			name := nameAny.(string)
			db.Name = &name
		}

		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (n *neo4jRedisRepo) GetNamespaces(ctx context.Context, dbUID string) ([]*types.RedisNamespace, error) {
	query := `MATCH
	(rn:RedisNS)-[:IN]->(rd:RedisDB)
	WHERE elementId(rd) = $uid
	RETURN rn
	`

	args := map[string]any{
		"uid": dbUID,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	var namespaces []*types.RedisNamespace

	if len(res.Records) == 0 {
		return namespaces, nil
	}

	for _, record := range res.Records {
		itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](record, "rn")
		if err != nil {
			return nil, err
		}

		namespace := &types.RedisNamespace{
			UID: &itemNode.ElementId,
		}

		if nameAny, ok := itemNode.Props["name"]; ok {
			name := nameAny.(string)
			namespace.Name = &name
		}

		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil
}

func (n *neo4jRedisRepo) Delete(ctx context.Context, uid string) error {
	query := `MATCH
	(r:Redis)
	WHERE elementId(r) = $uid
	DETACH DELETE r
	`

	args := map[string]any{
		"uid": uid,
	}

	_, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}

	return nil
}

func (n *neo4jRedisRepo) DeleteDB(ctx context.Context, uid string) error {
	query := `MATCH
	(rd:RedisDB)
	WHERE elementId(rd) = $uid
	DETACH DELETE rd
	`

	args := map[string]any{
		"uid": uid,
	}

	_, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}

	return nil
}

func (n *neo4jRedisRepo) DeleteNamespace(ctx context.Context, uid string) error {
	query := `MATCH
	(rn:RedisNS)
	WHERE elementId(rn) = $uid
	DETACH DELETE rn
	`

	args := map[string]any{
		"uid": uid,
	}

	_, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}

	return nil
}

func (n *neo4jRedisRepo) Update(ctx context.Context, redis *types.Redis) (*types.Redis, error) {
	query := `MATCH
	(r:Redis)
	WHERE r.host = $host
	SET
	`
	params := []string{}

	if redis.Host == nil {
		return nil, fmt.Errorf("redis cannot updated, host field is required")
	}
	if redis.Port != nil {
		params = append(params, "r.port = $port")
	}
	if redis.Master != nil {
		params = append(params, "r.master = $master")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN r"

	args := map[string]any{
		"host":   redis.Host,
		"port":   redis.Port,
		"master": redis.Master,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("redis not updated")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "r")
	if err != nil {
		return nil, err
	}

	newRedis := &types.Redis{
		UID: &itemNode.ElementId,
	}

	if hostAny, ok := itemNode.Props["host"]; ok {
		host := hostAny.(string)
		newRedis.Host = &host
	}
	if portAny, ok := itemNode.Props["port"]; ok {
		port := portAny.(int64)
		newRedis.Port = &port
	}
	if masterAny, ok := itemNode.Props["master"]; ok {
		master := masterAny.(string)
		newRedis.Master = &master
	}

	return newRedis, nil
}

func (n *neo4jRedisRepo) UpdateDB(ctx context.Context, redisDB *types.RedisDB) (*types.RedisDB, error) {
	query := `MATCH
	(rd:RedisDB)
	WHERE elementId(rd) = $uid
	SET
	`
	params := []string{}

	if redisDB.UID == nil {
		return nil, fmt.Errorf("redis db cannot be updated, uid field is required")
	}
	if redisDB.Name != nil {
		params = append(params, "rd.name = $name")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN rd"

	args := map[string]any{
		"uid":  redisDB.UID,
		"name": redisDB.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("redis db not updated")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "rd")
	if err != nil {
		return nil, err
	}

	newRedisDB := &types.RedisDB{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		newRedisDB.Name = &name
	}

	return newRedisDB, nil
}

func (n *neo4jRedisRepo) UpdateNamespace(ctx context.Context, redisNS *types.RedisNamespace) (*types.RedisNamespace, error) {
	query := `MATCH
	(rn:RedisNS)
	WHERE elementId(rn) = $uid
	SET
	`
	params := []string{}

	if redisNS.UID == nil {
		return nil, fmt.Errorf("redis ns cannot be updated, uid field is required")
	}
	if redisNS.Name != nil {
		params = append(params, "rn.name = $name")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN rn"

	args := map[string]any{
		"uid":  redisNS.UID,
		"name": redisNS.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("redis ns not updated")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "rn")
	if err != nil {
		return nil, err
	}

	newRedisNS := &types.RedisNamespace{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		newRedisNS.Name = &name
	}

	return newRedisNS, nil
}
