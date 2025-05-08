package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseSelect(stmt *pg_query.Node_SelectStmt, tables map[string]*types.Table) (*types.Select, error) {
	selectSt := &types.Select{
		Columns: []string{},
		Tables:  []string{},
	}

	for _, target := range stmt.SelectStmt.TargetList {
		switch targetNode := target.Node.(type) {
		case *pg_query.Node_ResTarget:
			switch valNode := targetNode.ResTarget.Val.Node.(type) {
			case *pg_query.Node_ColumnRef:
				for _, field := range valNode.ColumnRef.Fields {
					switch fieldNode := field.Node.(type) {
					case *pg_query.Node_String_:
						selectSt.Columns = append(selectSt.Columns, fieldNode.String_.Sval)
					case *pg_query.Node_AStar:
						selectSt.Columns = append(selectSt.Columns, "*")
					default:
						slog.Error("unimplemented column type", "type", fmt.Sprintf("%T", fieldNode))
					}
				}
			default:
				slog.Error("unimplemented column type", "type", fmt.Sprintf("%T", valNode))
			}
		default:
			slog.Error("unimplemented column type", "type", fmt.Sprintf("%T", targetNode))
		}
	}

	for _, fromClause := range stmt.SelectStmt.FromClause {
		switch node := fromClause.Node.(type) {
		case *pg_query.Node_RangeVar:
			selectSt.Tables = append(selectSt.Tables, node.RangeVar.Relname)
		default:
			slog.Error("unimplemented from clause type", "type", fmt.Sprintf("%T", node))
		}
	}

	return selectSt, nil
}
