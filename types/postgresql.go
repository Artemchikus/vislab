package types

type Postgresql struct {
	Host      *string
	Port      *int64
	Databases []*PostgresqlDB
	User      *string
}

type PostgresqlDB struct {
	Name          *string
	Owner         *string
	ForMigrations *bool
	Schemes       []*PostgresqlScheme
}

type PostgresqlScheme struct {
	Name   *string
	Owner  *string
	Tables []*PostgresqlTable
}

type PostgresqlTable struct {
	Name  *string
	Owner *string
	Type  *string
}
