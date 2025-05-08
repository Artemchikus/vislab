package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseCreateTrigger(stmt *pg_query.Node_CreateTrigStmt, triggers map[string]*types.Trigger) (*types.Trigger, error) {
	trigger := &types.Trigger{
		Name:  stmt.CreateTrigStmt.Trigname,
		Table: stmt.CreateTrigStmt.Relation.Relname,
	}

	for _, fName := range stmt.CreateTrigStmt.Funcname {
		switch fName := fName.Node.(type) {
		case *pg_query.Node_String_:
			trigger.Func = fName.String_.Sval
		default:
			slog.Error("unimplemented create trigger function name type", "type", fmt.Sprintf("%T", fName))
		}
	}

	return trigger, nil
}
