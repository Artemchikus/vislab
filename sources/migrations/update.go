package migrations

import (
	"fmt"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseUpdate(stmt *pg_query.Node_UpdateStmt, tables map[string]*types.Table) error {
	fmt.Printf("update_table: %s\n", stmt.UpdateStmt.Relation.Relname)

	return nil
}
