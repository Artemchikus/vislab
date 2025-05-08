package types

type RabbitMQ struct {
	Host   *string
	Port   *int64
	User   *string
	Queues []*RabbitQueue
}

type RabbitQueue struct {
	Name      *string
	QueueType *string
	Topic     *string
	TypeName  *string
}
