package migrations

import (
	"fmt"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseDelete(stmt *pg_query.Node_DeleteStmt, tables map[string]*types.Table) error {
	fmt.Printf("delete_from_table: %s\n", stmt.DeleteStmt.Relation.Relname)

	return nil
}
