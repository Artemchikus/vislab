package types

type Team struct {
	Name *string
	Devs []*Dev
}

type Dev struct {
	Name *string
	Link *string
	Role *string
}
