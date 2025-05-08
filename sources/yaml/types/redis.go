package types

type Redises struct {
	Instances    []*Redis `yaml:"instances"`
	LastInstance *Redis   `yaml:"-"`
}

type Redis struct {
	Host         *string    `yaml:"host"`
	Port         *int64     `yaml:"port"`
	Databases    []*RedisDB `yaml:"databases"`
	Master       *string    `yaml:"master"`
	Sentinel     *Sentinel  `yaml:"sentinel"`
	LastDatabase *RedisDB   `yaml:"-"`
}

type Sentinel struct {
	Host *string `yaml:"host"`
	Port *int64  `yaml:"port"`
}

type RedisDB struct {
	Name          *string           `yaml:"name"`
	Namespaces    []*RedisNamespace `yaml:"namespaces"`
	LastNamespace *RedisNamespace   `yaml:"-"`
}

type RedisNamespace struct {
	Name *string `yaml:"name"`
}
