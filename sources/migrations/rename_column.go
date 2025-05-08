package migrations

import (
	"fmt"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseRenameColumn(stmt *pg_query.Node_RenameStmt, tables map[string]*types.Table) (*types.Table, error) {
	table, ok := tables[stmt.RenameStmt.Relation.Relname]
	if !ok {
		return nil, fmt.Errorf("table %s not found", stmt.RenameStmt.Relation.Relname)
	}

	for _, column := range table.Columns {
		if column.Name == stmt.RenameStmt.Subname {
			column.Name = stmt.RenameStmt.Newname
		}
	}

	return table, nil
}
