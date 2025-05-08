package migrations

import (
	"fmt"
	"log/slog"
	"slices"
	"vislab/sources/migrations/types"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

func parseCreateTableAs(stmt *pg_query.Node_CreateTableAsStmt, tables map[string]*types.Table) (*types.Table, error) {
	table := &types.Table{
		Name:    stmt.CreateTableAsStmt.Into.Rel.Relname,
		Columns: []*types.Column{},
		Schema:  stmt.CreateTableAsStmt.Into.Rel.Schemaname,
		Type:    "common",
	}

	if table.Schema == "" {
		table.Schema = "public"
	}

	switch qNode := stmt.CreateTableAsStmt.Query.Node.(type) {
	case *pg_query.Node_SelectStmt:
		sel, err := parseSelect(qNode, tables)
		if err != nil {
			return nil, err
		}

		for _, sTable := range sel.Tables {
			tmpTable, ok := tables[sTable]
			if !ok {
				slog.Error("table not found in existing list", "table", table)
				continue
			}

			if sel.Columns[0] == "*" {
				table.Columns = append(table.Columns, tmpTable.Columns...)
				continue
			}

			for _, column := range sel.Columns {
				if slices.ContainsFunc(tmpTable.Columns, func(c *types.Column) bool {
					if c.Name == column {
						table.Columns = append(table.Columns, c)
						return true
					}

					return false
				}) {
				}
			}
		}
	default:
		slog.Error("unimplemented create table as query type", "type", fmt.Sprintf("%T", qNode))
	}

	return table, nil
}
