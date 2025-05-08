package types

import "vislab/libs/check"

const (
	RabbitMQClass    NodeClass = "RabbitMQ"
	RabbitQueueClass NodeClass = "RabbitQueue"
)

type RabbitMQ struct {
	UID  *string
	Host *string
	Port *int64
	User *string
}

func (r *RabbitMQ) Equal(other *RabbitMQ) bool {
	return check.ComparePointers(r.Host, other.Host) &&
		check.ComparePointers(r.Port, other.Port) &&
		check.ComparePointers(r.User, other.User)
}

type RabbitQueue struct {
	UID  *string
	Name *string
}

func (r *RabbitQueue) Equal(other *RabbitQueue) bool {
	return check.ComparePointers(r.Name, other.Name)
}
