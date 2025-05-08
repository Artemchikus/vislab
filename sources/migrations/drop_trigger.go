package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseDropTrigger(stmt *pg_query.Node_DropStmt, triggers map[string]*types.Trigger) error {
	for _, object := range stmt.DropStmt.Objects {
		switch oNode := object.Node.(type) {
		case *pg_query.Node_List:
			if len(oNode.List.Items) < 2 {
				slog.Error("expected trigger and table name", "items", len(oNode.List.Items))
				continue
			}

			switch iNode := oNode.List.Items[1].Node.(type) {
			case *pg_query.Node_String_:
				if _, ok := triggers[iNode.String_.Sval]; !ok {
					slog.Error("trigger not found in existing list", "trigger", iNode.String_.Sval)
					continue
				}

				delete(triggers, iNode.String_.Sval)
			default:
				slog.Error("unimplemented drop trigger item type", "type", fmt.Sprintf("%T", iNode))
			}
		default:
			slog.Error("unimplemented drop trigger type", "type", fmt.Sprintf("%T", oNode))
		}
	}
	return nil
}
