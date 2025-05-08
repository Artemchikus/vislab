package types

type Kafka struct {
	Name   *string
	Host   *string
	Port   *int64
	Queues []*KafkaQueue
}

type KafkaQueue struct {
	Name      *string
	QueueType *string
	Topic     *string
	TypeName  *string
}
