package migrations

import (
	"fmt"
	"log/slog"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseCreateTable(stmt *pg_query.Node_CreateStmt, tables map[string]*types.Table) (*types.Table, error) {
	table := &types.Table{
		Name:    stmt.CreateStmt.Relation.Relname,
		Schema:  stmt.CreateStmt.Relation.Schemaname,
		Columns: []*types.Column{},
		Type:    "common",
	}

	if table.Schema == "" {
		table.Schema = "public"
	}

	if stmt.CreateStmt.Partbound != nil {
		table.Type = "partition"
	}

	if stmt.CreateStmt.Partspec != nil {
		table.Type = "partitioned"
	}

	if stmt.CreateStmt.InhRelations != nil {
		for _, rel := range stmt.CreateStmt.InhRelations {
			switch rNode := rel.Node.(type) {
			case *pg_query.Node_RangeVar:
				tmpTable, ok := tables[rNode.RangeVar.Relname]
				if !ok {
					slog.Error("table not found in existing list", "table", table)
					continue
				}

				table.Columns = append(table.Columns, tmpTable.Columns...)
			default:
				slog.Error("unimplemented inherit relation type", "type", fmt.Sprintf("%T", rNode))
			}
		}

		return table, nil
	}

	for _, elt := range stmt.CreateStmt.TableElts {
		switch eNode := elt.Node.(type) {
		case *pg_query.Node_ColumnDef:
			column, err := parseColumn(eNode)
			if err != nil {
				slog.Error("failed to parse column", "error", err)
				continue
			}

			table.Columns = append(table.Columns, column)
		case *pg_query.Node_Constraint:
			switch eNode.Constraint.Contype {
			case pg_query.ConstrType_CONSTR_UNIQUE:
				for _, key := range eNode.Constraint.Keys {
					switch kNode := key.Node.(type) {
					case *pg_query.Node_String_:
						fmt.Printf("constraint_key: %v\n", kNode.String_.Sval)
					default:
						slog.Error("unimplemented constraint key type", "type", fmt.Sprintf("%T", kNode))
					}
				}
			default:
				slog.Error("unimplemented constraint type", "type", eNode.Constraint.Contype)
			}
		case *pg_query.Node_TableLikeClause:
			tmpTable, ok := tables[eNode.TableLikeClause.Relation.Relname]
			if !ok {
				slog.Error("table not found in existing list", "table", table)
				continue
			}

			table.Columns = append(table.Columns, tmpTable.Columns...)
		default:
			slog.Error("unimplemented elt type", "type", fmt.Sprintf("%T", eNode))
		}
	}

	return table, nil
}

func parseColumn(node *pg_query.Node_ColumnDef) (*types.Column, error) {
	column := &types.Column{
		Name: node.ColumnDef.Colname,
	}

	for _, typeName := range node.ColumnDef.TypeName.Names {
		switch tNode := typeName.Node.(type) {
		case *pg_query.Node_String_:
			if tNode.String_.Sval != "pg_catalog" {
				column.Type = tNode.String_.Sval
			}
		default:
			return nil, fmt.Errorf("unimplemented column_type type: %T", tNode)
		}
	}

	for _, constr := range node.ColumnDef.Constraints {
		switch cNode := constr.Node.(type) {
		case *pg_query.Node_Constraint:
			column.Constraints = append(column.Constraints, cNode.Constraint.Contype.String())
		default:
			slog.Error("unimplemented constraint type", "type", fmt.Sprintf("%T", cNode))
		}
	}

	return column, nil
}
