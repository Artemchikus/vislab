package migrations

import (
	"fmt"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseInsert(stmt *pg_query.Node_InsertStmt, tables map[string]*types.Table) error {
	fmt.Printf("insert_into_table: %s\n", stmt.InsertStmt.Relation.Relname)

	return nil
}
