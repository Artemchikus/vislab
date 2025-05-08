package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseAlterType(stmt *pg_query.Node_AlterEnumStmt, types map[string]*types.Type) (*types.Type, error) {
	for _, typeName := range stmt.AlterEnumStmt.TypeName {
		switch node := typeName.Node.(type) {
		case *pg_query.Node_String_:
			typ, ok := types[node.String_.Sval]
			if !ok {
				return nil, fmt.Errorf("type %s not found", node.String_.Sval)
			}
			delete(types, node.String_.Sval)

			typ.Name = node.String_.Sval

			return typ, nil
		default:
			slog.Error("unimplemented alter_type type", "type", fmt.Sprintf("%T", node))
		}
	}
	return nil, fmt.Errorf("type not found")
}
