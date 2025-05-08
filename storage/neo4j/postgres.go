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
	neo4jPostgresRepo struct {
		db neo4j.DriverWithContext
	}
)

func (n *Neo4jStorage) Postgres() storage.PostgresRepository {
	if n.postgresRepo != nil {
		return n.postgresRepo
	}

	n.postgresRepo = &neo4jPostgresRepo{db: n.db}
	return n.postgresRepo
}

func (n *neo4jPostgresRepo) Create(ctx context.Context, postgres *types.Postgresql) (string, error) {
	query := `CREATE
	(p:Postgres {
		host: $host,
		port: $port,
		user: $user
	})
	RETURN p
	`

	args := map[string]any{
		"host": postgres.Host,
		"port": postgres.Port,
		"user": postgres.User,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("postgres node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "p")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jPostgresRepo) CreateDB(ctx context.Context, db *types.PostgresqlDB) (string, error) {
	query := `CREATE
	(pd:PostgresDB {
		name: $name
	})
	RETURN pd
	`

	args := map[string]any{
		"name": db.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("postgres db node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "pd")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jPostgresRepo) CreateScheme(ctx context.Context, scheme *types.PostgresqlScheme) (string, error) {
	query := `CREATE
	(ps:PostgresScheme {
		name: $name
	})
	RETURN ps
	`

	args := map[string]any{
		"name": scheme.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("postgres scheme node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "ps")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jPostgresRepo) CreateTable(ctx context.Context, table *types.PostgresqlTable) (string, error) {
	query := `CREATE
	(pt:PostgresTable {
		name: $name
	})
	RETURN pt
	`

	args := map[string]any{
		"name": table.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return "", err
	}

	if len(res.Records) == 0 {
		return "", fmt.Errorf("postgres table node not created")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "pt")
	if err != nil {
		return "", err
	}

	return itemNode.ElementId, nil
}

func (n *neo4jPostgresRepo) Get(ctx context.Context, host string) (*types.Postgresql, error) {
	query := `MATCH
	(p:Postgres)
	WHERE p.host = $host
	RETURN p
	`

	args := map[string]any{
		"host": host,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("postgres node not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "p")
	if err != nil {
		return nil, err
	}

	postgres := &types.Postgresql{
		UID: &itemNode.ElementId,
	}

	if hostAny, ok := itemNode.Props["host"]; ok {
		host := hostAny.(string)
		postgres.Host = &host
	}
	if portAny, ok := itemNode.Props["port"]; ok {
		port := portAny.(int64)
		postgres.Port = &port
	}
	if userAny, ok := itemNode.Props["user"]; ok {
		user := userAny.(string)
		postgres.User = &user
	}

	return postgres, nil
}

func (n *neo4jPostgresRepo) GetDBs(ctx context.Context, postgresUID string) ([]*types.PostgresqlDB, error) {
	query := `MATCH
	(pd:PostgresDB)-[:IN]->(p:Postgres)
	WHERE elementId(p) = $uid
	RETURN pd
	`

	args := map[string]any{
		"uid": postgresUID,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	dbs := []*types.PostgresqlDB{}

	if len(res.Records) == 0 {
		return dbs, nil
	}

	for _, record := range res.Records {
		itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](record, "pd")
		if err != nil {
			return nil, err
		}

		db := &types.PostgresqlDB{
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

func (n *neo4jPostgresRepo) GetSchemes(ctx context.Context, dbUID string) ([]*types.PostgresqlScheme, error) {
	query := `MATCH
	(ps:PostgresScheme)-[:IN]->(pd:PostgresDB)
	WHERE elementId(pd) = $uid
	RETURN ps
	`

	args := map[string]any{
		"uid": dbUID,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	schemes := []*types.PostgresqlScheme{}

	if len(res.Records) == 0 {
		return schemes, nil
	}

	for _, record := range res.Records {
		itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](record, "ps")
		if err != nil {
			return nil, err
		}

		scheme := &types.PostgresqlScheme{
			UID: &itemNode.ElementId,
		}

		if nameAny, ok := itemNode.Props["name"]; ok {
			name := nameAny.(string)
			scheme.Name = &name
		}

		schemes = append(schemes, scheme)
	}

	return schemes, nil
}

func (n *neo4jPostgresRepo) GetTables(ctx context.Context, schemeUID string) ([]*types.PostgresqlTable, error) {
	query := `MATCH
	(pt:PostgresTable)-[:IN]->(ps:PostgresScheme)
	WHERE elementId(ps) = $uid
	RETURN pt
	`

	args := map[string]any{
		"uid": schemeUID,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	tables := []*types.PostgresqlTable{}

	if len(res.Records) == 0 {
		return tables, nil
	}

	for _, record := range res.Records {
		itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](record, "pt")
		if err != nil {
			return nil, err
		}

		table := &types.PostgresqlTable{
			UID: &itemNode.ElementId,
		}

		if nameAny, ok := itemNode.Props["name"]; ok {
			name := nameAny.(string)
			table.Name = &name
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func (n *neo4jPostgresRepo) Delete(ctx context.Context, uid string) error {
	query := `MATCH
	(p:Postgres)
	WHERE elementId(p) = $uid
	DETACH DELETE p
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

func (n *neo4jPostgresRepo) DeleteDB(ctx context.Context, uid string) error {
	query := `MATCH
	(pd:PostgresDB)
	WHERE elementId(pd) = $uid
	DETACH DELETE pd
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

func (n *neo4jPostgresRepo) DeleteScheme(ctx context.Context, uid string) error {
	query := `MATCH
	(ps:PostgresScheme)
	WHERE elementId(ps) = $uid
	DETACH DELETE ps
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

func (n *neo4jPostgresRepo) DeleteTable(ctx context.Context, uid string) error {
	query := `MATCH
	(pt:PostgresTable)
	WHERE elementId(pt) = $uid
	DETACH DELETE pt
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

func (n *neo4jPostgresRepo) Update(ctx context.Context, postgres *types.Postgresql) (*types.Postgresql, error) {
	query := `MATCH
	(p:Postgres)
	WHERE p.host = $host
	SET
	`

	params := []string{}

	if postgres.Host == nil {
		return nil, fmt.Errorf("postgres cannot be updated, host field is required")
	}
	if postgres.Port != nil {
		params = append(params, "p.port = $port")
	}
	if postgres.User != nil {
		params = append(params, "p.user = $user")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN p"

	args := map[string]any{
		"host": postgres.Host,
		"port": postgres.Port,
		"user": postgres.User,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("postgres node not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "p")
	if err != nil {
		return nil, err
	}

	newPostgres := &types.Postgresql{
		UID: &itemNode.ElementId,
	}

	if hostAny, ok := itemNode.Props["host"]; ok {
		host := hostAny.(string)
		newPostgres.Host = &host
	}
	if portAny, ok := itemNode.Props["port"]; ok {
		port := portAny.(int64)
		newPostgres.Port = &port
	}
	if userAny, ok := itemNode.Props["user"]; ok {
		user := userAny.(string)
		newPostgres.User = &user
	}

	return newPostgres, nil
}

func (n *neo4jPostgresRepo) UpdateDB(ctx context.Context, db *types.PostgresqlDB) (*types.PostgresqlDB, error) {
	query := `MATCH
	(pd:PostgresDB)
	WHERE elementId(pd) = $uid
	SET
	`

	params := []string{}

	if db.UID == nil {
		return nil, fmt.Errorf("postgres db cannot be updated, uid field is required")
	}
	if db.Name != nil {
		params = append(params, "pd.name = $name")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("noting to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN pd"

	args := map[string]any{
		"uid":  db.UID,
		"name": db.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("postgres db node not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "pd")
	if err != nil {
		return nil, err
	}

	newDB := &types.PostgresqlDB{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		newDB.Name = &name
	}

	return newDB, nil
}

func (n *neo4jPostgresRepo) UpdateScheme(ctx context.Context, scheme *types.PostgresqlScheme) (*types.PostgresqlScheme, error) {
	query := `MATCH
	(ps:PostgresScheme)
	WHERE elementId(ps) = $uid
	SET
	`

	params := []string{}

	if scheme.UID == nil {
		return nil, fmt.Errorf("postgres scheme cannot be updated, uid field is required")
	}
	if scheme.Name != nil {
		params = append(params, "ps.name = $name")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("noting to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN ps"

	args := map[string]any{
		"uid":  scheme.UID,
		"name": scheme.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("postgres scheme node not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "ps")
	if err != nil {
		return nil, err
	}

	newScheme := &types.PostgresqlScheme{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		newScheme.Name = &name
	}

	return newScheme, nil
}

func (n *neo4jPostgresRepo) UpdateTable(ctx context.Context, table *types.PostgresqlTable) (*types.PostgresqlTable, error) {
	query := `MATCH
	(pt:PostgresTable)
	WHERE elementId(pt) = $uid
	SET
	`

	params := []string{}

	if table.Name == nil {
		return nil, fmt.Errorf("postgres table cannot be updated, name field is required")
	}

	if len(params) == 0 {
		return nil, fmt.Errorf("noting to update")
	}

	query += strings.Join(params, ", ")
	query += " RETURN pt"

	args := map[string]any{
		"uid": table.Name,
	}

	res, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	if len(res.Records) == 0 {
		return nil, fmt.Errorf("postgres table node not found")
	}

	itemNode, _, err := neo4j.GetRecordValue[neo4j.Node](res.Records[0], "pt")
	if err != nil {
		return nil, err
	}

	newTable := &types.PostgresqlTable{
		UID: &itemNode.ElementId,
	}

	if nameAny, ok := itemNode.Props["name"]; ok {
		name := nameAny.(string)
		newTable.Name = &name
	}

	return newTable, nil
}
