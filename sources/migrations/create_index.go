package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseCreateIndex(stmt *pg_query.Node_IndexStmt, indexes map[string]*types.Index) (*types.Index, error) {
	index := &types.Index{
		Name:  stmt.IndexStmt.Idxname,
		Table: stmt.IndexStmt.Relation.Relname,
	}

	for _, column := range stmt.IndexStmt.IndexParams {
		switch node := column.Node.(type) {
		case *pg_query.Node_IndexElem:
			index.Columns = append(index.Columns, node.IndexElem.Name)
		default:
			slog.Error("unimplemented index column type", "type", fmt.Sprintf("%T", node))
		}
	}

	return index, nil
}
