package types

import "vislab/libs/check"

const (
	RedisClass   NodeClass = "Redis"
	RedisDBClass NodeClass = "RedisDB"
	RedisNSClass NodeClass = "RedisNS"
)

type Redis struct {
	UID      *string
	Host     *string
	Port     *int64
	Master   *string
	Sentinel *Sentinel // ignoring for now
}

func (r *Redis) Equal(other *Redis) bool {
	return check.ComparePointers(r.Host, other.Host) &&
		check.ComparePointers(r.Port, other.Port) &&
		check.ComparePointers(r.Master, other.Master)
}

type Sentinel struct {
	Host string
	Port int64
}

func (s *Sentinel) Equal(other *Sentinel) bool {
	return s.Host == other.Host &&
		s.Port == other.Port
}

type RedisDB struct {
	UID  *string
	Name *string
}

func (r *RedisDB) Equal(other *RedisDB) bool {
	return check.ComparePointers(r.Name, other.Name)
}

type RedisNamespace struct {
	UID  *string
	Name *string
}

func (r *RedisNamespace) Equal(other *RedisNamespace) bool {
	return check.ComparePointers(r.Name, other.Name)
}
