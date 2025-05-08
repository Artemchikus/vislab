package types

type RabbitMQs struct {
	Instances    []*RabbitMQ `yaml:"instances"`
	LastInstance *RabbitMQ   `yaml:"-"`
}

type RabbitMQ struct {
	Host      *string        `yaml:"host"`
	Port      *int64         `yaml:"port"`
	User      *string        `yaml:"user"`
	Queues    []*RabbitQueue `yaml:"queues"`
	LastQueue *RabbitQueue   `yaml:"-"`
}

type RabbitQueue struct {
	Name *string `yaml:"name"`
}
