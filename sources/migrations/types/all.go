package types

type (
	All struct {
		Tables   map[string]*Table
		Funcs    map[string][]*Func
		Indexes  map[string]*Index
		Triggers map[string]*Trigger
		Types    map[string]*Type
	}
)
