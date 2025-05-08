package types

import "vislab/libs/check"

const (
	KafkaClass      NodeClass = "Kafka"
	KafkaQueueClass NodeClass = "KafkaQueue"
)

type Kafka struct {
	UID  *string
	Name *string
	Host *string
	Port *int64
}

func (k *Kafka) Equal(other *Kafka) bool {
	return check.ComparePointers(k.Name, other.Name) &&
		check.ComparePointers(k.Host, other.Host) &&
		check.ComparePointers(k.Port, other.Port)
}

type KafkaQueue struct {
	UID       *string
	Name      *string
	QueueType *string
	Topic     *string
	TypeName  *string
}

func (kq *KafkaQueue) Equal(other *KafkaQueue) bool {
	return check.ComparePointers(kq.Name, other.Name) &&
		check.ComparePointers(kq.QueueType, other.QueueType) &&
		check.ComparePointers(kq.Topic, other.Topic) &&
		check.ComparePointers(kq.TypeName, other.TypeName)
}
