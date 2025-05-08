package types

type Redis struct {
	Host      *string
	Port      *int64
	Databases []*RedisDB
	Master    *string
	Sentinel  *Sentinel
}

type Sentinel struct {
	Host *string
	Port *int64
}

type RedisDB struct {
	Name       *string
	Owner      *string
	Namespaces []*RedisNamespace
}

type RedisNamespace struct {
	Name  *string
	Owner *string
}
