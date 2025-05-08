package types

type (
	Table struct {
		Name    string
		Schema  string
		Columns []*Column
		Type    string
	}
)
