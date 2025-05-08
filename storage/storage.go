package storage

import "context"

type Storage interface {
	// Reconnect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Service() ServiceRepository
	// Team() TeamRepository
	// Pipeline() PipelineRepository
	Kafka() KafkaRepository
	Redis() RedisRepository
	RabbitMQ() RabbitMQRepository
	Postgres() PostgresRepository
	Connection() ConnectionRepository
}
