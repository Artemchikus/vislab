package migrations

import (
	"fmt"
	"log/slog"
	"slices"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseAlterTable(stmt *pg_query.Node_AlterTableStmt, tables map[string]*types.Table) (*types.Table, error) {
	table, ok := tables[stmt.AlterTableStmt.Relation.Relname]
	if !ok {
		slog.Error("table not found in existing list", "table", stmt.AlterTableStmt.Relation.Relname)
		return nil, fmt.Errorf("table not found in existing list: %s", stmt.AlterTableStmt.Relation.Relname)
	}

	for _, cmd := range stmt.AlterTableStmt.Cmds {
		switch cNode := cmd.Node.(type) {
		case *pg_query.Node_AlterTableCmd:
			switch cNode.AlterTableCmd.Subtype {
			case pg_query.AlterTableType_AT_AddConstraint:
				switch def := cNode.AlterTableCmd.Def.Node.(type) {
				case *pg_query.Node_Constraint:
					for _, key := range def.Constraint.Keys {
						switch kNode := key.Node.(type) {
						case *pg_query.Node_String_:
							for _, column := range table.Columns {
								if column.Name == kNode.String_.Sval {
									column.Constraints = append(column.Constraints, def.Constraint.Contype.String())
								}
							}
						default:
							slog.Error("unimplemented constraint key type", "type", fmt.Sprintf("%T", kNode))
						}
					}
				default:
					slog.Error("unimplemented constraint type", "type", fmt.Sprintf("%T", def))
				}
			case pg_query.AlterTableType_AT_AddColumn:
				switch def := cNode.AlterTableCmd.Def.Node.(type) {
				case *pg_query.Node_ColumnDef:
					column, err := parseColumn(def)
					if err != nil {
						return nil, err
					}

					if slices.ContainsFunc(table.Columns, func(c *types.Column) bool {
						return c.Name == column.Name
					}) {
						return nil, fmt.Errorf("column already exists")
					}

					table.Columns = append(table.Columns, column)
				default:
					slog.Error("unimplemented add column type", "type", fmt.Sprintf("%T", def))
				}
			case pg_query.AlterTableType_AT_AlterColumnType:
				switch def := cNode.AlterTableCmd.Def.Node.(type) {
				case *pg_query.Node_ColumnDef:
					column, err := parseColumn(def)
					if err != nil {
						return nil, err
					}

					for _, v := range table.Columns {
						if v.Name == column.Name {
							v.Name = cNode.AlterTableCmd.Name
						}
					}
				default:
					slog.Error("unimplemented alter column type", "type", fmt.Sprintf("%T", def))
				}
			case pg_query.AlterTableType_AT_DropConstraint:
				fmt.Printf("drop_constraint: %v\n", cNode.AlterTableCmd.Name)
			case pg_query.AlterTableType_AT_DropColumn:
				for i, column := range table.Columns {
					if column.Name == cNode.AlterTableCmd.Name {
						table.Columns = append(table.Columns[:i], table.Columns[i+1:]...)
					}
				}
			case pg_query.AlterTableType_AT_ColumnDefault:
				fmt.Printf("field default set, skipping: %v\n", cNode) // TODO add set default handler
			case pg_query.AlterTableType_AT_DropNotNull:
				fmt.Printf("not null field dropping, skipping: %v\n", cNode)
			case pg_query.AlterTableType_AT_AttachPartition:
				switch def := cNode.AlterTableCmd.Def.Node.(type) {
				case *pg_query.Node_PartitionCmd:
					table.Type = "partitioned"
					partTable, ok := tables[def.PartitionCmd.Name.Relname]
					if !ok {
						slog.Error("table not found in existing list", "table", def.PartitionCmd.Name.Relname)
						return nil, fmt.Errorf("table not found in existing list: %s", def.PartitionCmd.Name.Relname)
					}

					partTable.Type = "partition"
				default:
					slog.Error("unimplemented attach partition type", "type", fmt.Sprintf("%T", def))
				}
			default:
				slog.Error("unimplemented alter table type", "type", cNode.AlterTableCmd.Subtype)
			}
		default:
			slog.Error("unimplemented alter table command type", "type", fmt.Sprintf("%T", cNode))
		}
	}

	return table, nil
}
