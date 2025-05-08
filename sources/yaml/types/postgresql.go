package types

type Postgresqls struct {
	Instances    []*Postgresql `yaml:"instances"`
	LastInstance *Postgresql   `yaml:"-"`
}

type Postgresql struct {
	Host         *string `yaml:"host"`
	Port         *int64  `yaml:"port"`
	User         *string `yaml:"user"`
	Databases    []*PqDB `yaml:"databases"`
	LastDatabase *PqDB   `yaml:"-"`
}

type PqDB struct {
	Name          *string     `yaml:"name"`
	ForMigrations *bool       `yaml:"for_migrations"`
	Schemes       []*PqScheme `yaml:"schemes"`
	LastScheme    *PqScheme   `yaml:"-"`
}

type PqScheme struct {
	Name *string `yaml:"name"`
}
