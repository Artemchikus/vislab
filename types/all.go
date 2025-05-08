package types

type (
	All struct {
		Service       *Service
		Kafkas        []*Kafka
		Redises       []*Redis
		Postgresqls   []*Postgresql
		RabbitMQs     []*RabbitMQ
		OtherServices []*Service
	}
)