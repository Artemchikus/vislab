package types

type (
	ConnType  string
	NodeClass string
	ConnNode  struct {
		Class NodeClass
		ID    string
	}
)

const (
	ConnIN           ConnType = "IN"
	ConnSendsTo      ConnType = "SENDS_TO"
	ConnReceivesFrom ConnType = "RECEIVES_FROM"
	ConnUses         ConnType = "USES"
)

func (c ConnType) String() string {
	return string(c)
}

func (n NodeClass) String() string {
	return string(n)
}

func (c ConnNode) Equal(other *ConnNode) bool {
	return c.Class == other.Class
}
