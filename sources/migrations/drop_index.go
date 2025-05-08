package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseDropIndex(stmt *pg_query.Node_DropStmt, indexes map[string]*types.Index) error {
	for _, object := range stmt.DropStmt.Objects {
		switch oNode := object.Node.(type) {
		case *pg_query.Node_List:
			for _, item := range oNode.List.Items {
				switch iNode := item.Node.(type) {
				case *pg_query.Node_String_:
					if _, ok := indexes[iNode.String_.Sval]; !ok {
						slog.Error("index not found in existing list", "index", iNode.String_.Sval)
						continue
					}

					delete(indexes, iNode.String_.Sval)
				default:
					slog.Error("unimplemented drop index item type", "type", fmt.Sprintf("%T", iNode))
				}
			}
		default:
			slog.Error("unimplemented drop index type", "type", fmt.Sprintf("%T", oNode))
		}
	}
	return nil
}
