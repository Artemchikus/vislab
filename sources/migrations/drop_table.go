package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseDropTable(stmt *pg_query.Node_DropStmt, tables map[string]*types.Table, indexes map[string]*types.Index, triggers map[string]*types.Trigger) error {
	for _, object := range stmt.DropStmt.Objects {
		switch oNode := object.Node.(type) {
		case *pg_query.Node_List:
			for _, item := range oNode.List.Items {
				switch iNode := item.Node.(type) {
				case *pg_query.Node_String_:
					if _, ok := tables[iNode.String_.Sval]; !ok {
						slog.Error("table not found in existing list", "table", iNode.String_.Sval)
						continue
					}

					delete(tables, iNode.String_.Sval)

					for _, index := range indexes {
						if index.Table == iNode.String_.Sval {
							delete(indexes, index.Name)
						}
					}

					for _, trigger := range triggers {
						if trigger.Table == iNode.String_.Sval {
							delete(triggers, trigger.Name)
						}
					}
				default:
					slog.Error("unimplemented drop table item type", "type", fmt.Sprintf("%T", iNode))
				}
			}
		default:
			slog.Error("unimplemented drop table type", "type", fmt.Sprintf("%T", oNode))
		}
	}
	return nil
}
