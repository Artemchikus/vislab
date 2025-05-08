package types

type (
	Trigger struct {
		Name    string
		Table   string
		Columns []string
		Func    string
	}
)
