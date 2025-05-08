package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseCreateFunc(stmt *pg_query.Node_CreateFunctionStmt, funcs map[string][]*types.Func) (*types.Func, error) {
	function := &types.Func{}

	for _, funcName := range stmt.CreateFunctionStmt.Funcname {
		switch fNode := funcName.Node.(type) {
		case *pg_query.Node_String_:
			function.Name = fNode.String_.Sval
		default:
			slog.Error("unimplemented create function name type", "type", fmt.Sprintf("%T", fNode))
		}
	}

	for _, param := range stmt.CreateFunctionStmt.Parameters {
		switch pNode := param.Node.(type) {
		case *pg_query.Node_FunctionParameter:
			for _, name := range pNode.FunctionParameter.ArgType.Names {
				switch nNode := name.Node.(type) {
				case *pg_query.Node_String_:
					if nNode.String_.Sval != "pg_catalog" {
						function.Args = append(function.Args, nNode.String_.Sval)
					}
				default:
					slog.Error("unimplemented create function parameter type", "type", fmt.Sprintf("%T", nNode))
				}
			}
		default:
			slog.Error("unimplemented create function parameter type", "type", fmt.Sprintf("%T", pNode))
		}
	}

	return function, nil
}
