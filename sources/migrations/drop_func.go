package migrations

import (
	"fmt"
	"log/slog"
	"slices"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseDropFunc(stmt *pg_query.Node_DropStmt, funcs map[string][]*types.Func) error {
	for _, object := range stmt.DropStmt.Objects {
		switch oNode := object.Node.(type) {
		case *pg_query.Node_ObjectWithArgs:
			function, err := parseFunc(oNode, funcs)
			if err != nil {
				slog.Error("failed to parse function", "error", err)
				continue
			}

			fs, ok := funcs[function.Name]
			if !ok {
				slog.Error("function not found in existing list", "function", function.Name)
				continue
			}

			for i, f := range fs {
				if f.Equal(function) {
					funcs[function.Name] = slices.Delete(fs, i, i)
					break
				}
			}
		default:
			slog.Error("unimplemented drop function type", "type", fmt.Sprintf("%T", oNode))
		}
	}

	return nil
}

func parseFunc(stmt *pg_query.Node_ObjectWithArgs, funcs map[string][]*types.Func) (*types.Func, error) {
	function := &types.Func{
		Args: []string{},
	}

	for _, objName := range stmt.ObjectWithArgs.Objname {
		switch oNode := objName.Node.(type) {
		case *pg_query.Node_String_:
			function.Name = oNode.String_.Sval
		default:
			return nil, fmt.Errorf("unimplemented drop function type: %T", oNode)
		}
	}

	if stmt.ObjectWithArgs.ArgsUnspecified {
		return function, nil
	}

	for _, objArg := range stmt.ObjectWithArgs.Objargs {
		switch oNode := objArg.Node.(type) {
		case *pg_query.Node_TypeName:
			for _, name := range oNode.TypeName.Names {
				switch nNode := name.Node.(type) {
				case *pg_query.Node_String_:
					if nNode.String_.Sval != "pg_catalog" {
						function.Args = append(function.Args, nNode.String_.Sval)
					}
				default:
					slog.Error("unimplemented drop function item type", "type", fmt.Sprintf("%T", nNode))
				}
			}
		default:
			slog.Error("unimplemented drop function type", "type", fmt.Sprintf("%T", oNode))
		}
	}

	return function, nil
}
