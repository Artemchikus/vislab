package migrations

import (
	"fmt"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseRenameIndex(stmt *pg_query.Node_RenameStmt, indexes map[string]*types.Index) (*types.Index, error) {
	index, ok := indexes[stmt.RenameStmt.Relation.Relname]
	if !ok {
		return nil, fmt.Errorf("index %s not found", stmt.RenameStmt.Relation.Relname)
	}
	delete(indexes, stmt.RenameStmt.Relation.Relname)

	index.Name = stmt.RenameStmt.Newname

	return index, nil
}
