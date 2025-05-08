package types

import "slices"

type (
	Func struct {
		Name string
		Args []string
	}
)

func (f *Func) Equal(other *Func) bool {
	return f.Name == other.Name && slices.Equal(f.Args, other.Args)
}
