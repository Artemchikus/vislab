package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseDropType(stmt *pg_query.Node_DropStmt, typs map[string]*types.Type) error {
	for _, object := range stmt.DropStmt.Objects {
		switch oNode := object.Node.(type) {
		case *pg_query.Node_TypeName:
			for _, name := range oNode.TypeName.Names {
				switch nNode := name.Node.(type) {
				case *pg_query.Node_String_:
					if _, ok := typs[nNode.String_.Sval]; !ok {
						slog.Error("type not found in existing list", "type", nNode.String_.Sval)
						continue
					}

					delete(typs, nNode.String_.Sval)
				default:
					slog.Error("unimplemented drop_type item type", "type", fmt.Sprintf("%T", nNode))
				}
			}
		default:
			slog.Error("unimplemented drop_type type", "type", fmt.Sprintf("%T", oNode))
		}
	}
	return nil
}
