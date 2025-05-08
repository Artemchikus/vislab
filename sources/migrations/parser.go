package migrations

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

type Parser struct {
}

func NewParser() (*Parser, error) {
	p := &Parser{}

	return p, nil
}

func (p *Parser) Parse(in []byte, out *types.All) error {
	str := strings.ReplaceAll(string(in), "-- SQL in this section is executed when the migration is rolled back.\n", "")
	str = strings.ReplaceAll(str, "-- SQL in this section is executed when the migration is applied.\n", "")
	str = strings.ReplaceAll(str, "-- +goose StatementBegin\n", "")
	str = strings.ReplaceAll(str, "-- +goose StatementEnd\n", "")

	str = strings.TrimPrefix(str, "-- +goose Up\n")
	parts := strings.SplitN(str, "-- +goose Down\n", 2)

	if len(parts) != 2 {
		return fmt.Errorf("file does not contain up and down migrations")
	}

	stmt, err := pg_query.Parse(parts[0])
	if err != nil {
		return err
	}

	for _, stmt := range stmt.Stmts {
		switch stmt := stmt.Stmt.Node.(type) {
		case *pg_query.Node_CreateStmt:
			table, err := parseCreateTable(stmt, out.Tables)
			if err != nil {
				return err
			}

			out.Tables[table.Name] = table

		case *pg_query.Node_DropStmt:
			switch stmt.DropStmt.RemoveType {
			case pg_query.ObjectType_OBJECT_TABLE:
				if err := parseDropTable(stmt, out.Tables, out.Indexes, out.Triggers); err != nil {
					return err
				}
			case pg_query.ObjectType_OBJECT_INDEX:
				if err := parseDropIndex(stmt, out.Indexes); err != nil {
					return err
				}
			case pg_query.ObjectType_OBJECT_FUNCTION:
				if err := parseDropFunc(stmt, out.Funcs); err != nil {
					return err
				}
			case pg_query.ObjectType_OBJECT_TRIGGER:
				if err := parseDropTrigger(stmt, out.Triggers); err != nil {
					return err
				} // TODO: add support for triggers ("on table" reads as another trigger)
			case pg_query.ObjectType_OBJECT_TYPE:
				if err := parseDropType(stmt, out.Types); err != nil {
					return err
				}
			default:
				slog.Error("unimplemented drop type", "type", stmt.DropStmt.RemoveType)
			}
		case *pg_query.Node_CreateFunctionStmt:
			function, err := parseCreateFunc(stmt, out.Funcs)
			if err != nil {
				return err
			}

			fs, ok := out.Funcs[function.Name]
			if !ok {
				out.Funcs[function.Name] = []*types.Func{function}
			} else {
				if !slices.ContainsFunc(fs, func(f *types.Func) bool {
					return f.Equal(function)
				}) {
					out.Funcs[function.Name] = append(fs, function)
				}
			}
		case *pg_query.Node_IndexStmt:
			index, err := parseCreateIndex(stmt, out.Indexes)
			if err != nil {
				return err
			}

			out.Indexes[index.Name] = index
		case *pg_query.Node_AlterTableStmt:
			table, err := parseAlterTable(stmt, out.Tables)
			if err != nil {
				return err
			}

			fmt.Printf("alter_table: %v\n", table.Name)
		case *pg_query.Node_CreateTrigStmt:
			trigger, err := parseCreateTrigger(stmt, out.Triggers)
			if err != nil {
				return err
			}

			out.Triggers[trigger.Name] = trigger
		case *pg_query.Node_CommentStmt:
			var commentType string
			switch stmt.CommentStmt.Objtype {
			case pg_query.ObjectType_OBJECT_TABLE:
				commentType = "table"
			case pg_query.ObjectType_OBJECT_INDEX:
				commentType = "index"
			case pg_query.ObjectType_OBJECT_COLUMN:
				commentType = "column"
			default:
				slog.Error("unimplemented comment type", "type", commentType)
			}
			switch node := stmt.CommentStmt.Object.Node.(type) {
			case *pg_query.Node_List:
				objects := []string{}

				for _, item := range node.List.Items {
					switch itemNode := item.Node.(type) {
					case *pg_query.Node_String_:
						objects = append(objects, itemNode.String_.Sval)
					default:
						slog.Error("unimplemented comment item type", "type", fmt.Sprintf("%T", itemNode))
					}
				}
			default:
				slog.Error("unimplemented comment object type", "type", fmt.Sprintf("%T", node))
			}
		case *pg_query.Node_RenameStmt:
			switch stmt.RenameStmt.RenameType {
			case pg_query.ObjectType_OBJECT_TABLE:
				table, err := parseRenameTable(stmt, out.Tables)
				if err != nil {
					return err
				}

				out.Tables[table.Name] = table
			case pg_query.ObjectType_OBJECT_INDEX:
				index, err := parseRenameIndex(stmt, out.Indexes)
				if err != nil {
					return err
				}

				out.Indexes[index.Name] = index
			case pg_query.ObjectType_OBJECT_COLUMN:
				column, err := parseRenameColumn(stmt, out.Tables)
				if err != nil {
					return err
				}

				fmt.Printf("rename_column_in_table: %v\n", column)
			default:
				slog.Error("unimplemented rename type", "type", stmt.RenameStmt.RenameType)
			}
		case *pg_query.Node_UpdateStmt:
			if err := parseUpdate(stmt, out.Tables); err != nil {
				return err
			}
		case *pg_query.Node_InsertStmt:
			if err := parseInsert(stmt, out.Tables); err != nil {
				return err
			}
		case *pg_query.Node_DeleteStmt:
			if err := parseDelete(stmt, out.Tables); err != nil {
				return err
			}
		case *pg_query.Node_CreateEnumStmt:
			ty, err := parseCreateType(stmt, out.Types)
			if err != nil {
				return err
			}

			out.Types[ty.Name] = ty
		case *pg_query.Node_AlterEnumStmt:
			ty, err := parseAlterType(stmt, out.Types)
			if err != nil {
				return err
			}

			out.Types[ty.Name] = ty
		case *pg_query.Node_CreateTableAsStmt:
			table, err := parseCreateTableAs(stmt, out.Tables)
			if err != nil {
				return err
			}

			out.Tables[table.Name] = table
		default:
			slog.Error("unimplemented type", "type", fmt.Sprintf("%T", stmt))
		}
	}

	return nil
}
