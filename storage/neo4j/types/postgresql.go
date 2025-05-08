package types

import "vislab/libs/check"

const (
	PostgresClass       NodeClass = "Postgres"
	PostgresDBClass     NodeClass = "PostgresDB"
	PostgresSchemeClass NodeClass = "PostgresScheme"
	PostgresTableClass  NodeClass = "PostgresTable"
)

type Postgresql struct {
	UID  *string
	Host *string
	Port *int64
	User *string
}

func (p *Postgresql) Equal(other *Postgresql) bool {
	return check.ComparePointers(p.Host, other.Host) &&
		check.ComparePointers(p.Port, other.Port) &&
		check.ComparePointers(p.User, other.User)
}

type PostgresqlDB struct {
	UID  *string
	Name *string
}

func (p *PostgresqlDB) Equal(other *PostgresqlDB) bool {
	return check.ComparePointers(p.Name, other.Name)
}

type PostgresqlScheme struct {
	UID  *string
	Name *string
}

func (p *PostgresqlScheme) Equal(other *PostgresqlScheme) bool {
	return check.ComparePointers(p.Name, other.Name)
}

type PostgresqlTable struct {
	UID  *string
	Name *string
}

func (p *PostgresqlTable) Equal(other *PostgresqlTable) bool {
	return check.ComparePointers(p.Name, other.Name)
}
