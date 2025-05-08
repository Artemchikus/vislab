package types

type Kafkas struct {
	Instances    []*Kafka `yaml:"instances"`
	LastInstance *Kafka   `yaml:"-"`
}

type Kafka struct {
	Name      *string       `yaml:"name"`
	Host      *string       `yaml:"host"`
	Port      *int64        `yaml:"port"`
	Queues    []*KafkaQueue `yaml:"queues"`
	LastQueue *KafkaQueue   `yaml:"-"`
}

type KafkaQueue struct {
	Name      *string `yaml:"name"`
	QueueType *string `yaml:"queue_type"`
	Topic     *string `yaml:"topic"`
	TypeName  *string `yaml:"type_name"`
}
