package types

type All struct {
	Service      *Services      `yaml:"service"`
	OtherService *OtherServices `yaml:"other_service"`
	Kafka        *Kafkas        `yaml:"kafka"`
	Redis        *Redises       `yaml:"redis"`
	Postgresql   *Postgresqls   `yaml:"postgresql"`
	RabbitMQ     *RabbitMQs     `yaml:"rabbitmq"`
}
