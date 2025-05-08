package neo4j

import (
	"context"
	"fmt"
	"vislab/storage"
	"vislab/storage/neo4j/types"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type (
	neo4jConnRepo struct {
		db neo4j.DriverWithContext
	}
)

func (n *Neo4jStorage) Connection() storage.ConnectionRepository {
	if n.connRepo != nil {
		return n.connRepo
	}

	n.connRepo = &neo4jConnRepo{db: n.db}
	return n.connRepo
}

func (n *neo4jConnRepo) Create(ctx context.Context, fromID, toID *types.ConnNode, connType types.ConnType) error {
	query := fmt.Sprintf(`MATCH
	(n:%s),
	(m:%s)
	WHERE elementId(n) = $fromID and elementId(m) = $toID
	MERGE
	(n)-[:%s]-(m)
	`, fromID.Class, toID.Class, connType.String())

	args := map[string]any{
		"fromID": fromID.ID,
		"toID":   toID.ID,
	}

	_, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}

	return nil
}

func (n *neo4jConnRepo) Delete(ctx context.Context, fromID, toID *types.ConnNode, connType types.ConnType) error {
	query := fmt.Sprintf(`MATCH
	(n:%s)-[c:%s]-(m:%s)
	WHERE elementId(n) = $fromID and elementId(m) = $toID
	DELETE c
	`, fromID.Class, toID.Class, connType.String())

	args := map[string]any{
		"fromID": fromID.ID,
		"toID":   toID.ID,
	}

	_, err := neo4j.ExecuteQuery(ctx, n.db, query, args, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}

	return nil
}
