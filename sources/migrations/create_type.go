package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseCreateType(stmt *pg_query.Node_CreateEnumStmt, typs map[string]*types.Type) (*types.Type, error) {
	typ := &types.Type{}

	for _, typeName := range stmt.CreateEnumStmt.TypeName {
		switch node := typeName.Node.(type) {
		case *pg_query.Node_String_:
			typ.Name = node.String_.Sval
		default:
			slog.Error("unimplemented create_type type", "type", fmt.Sprintf("%T", node))
		}
	}

	return typ, nil
}
