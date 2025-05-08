package storage

import (
	"context"
	"vislab/storage/neo4j/types"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *types.Service) (string, error)
	Get(ctx context.Context, name string) (*types.Service, error)
	Delete(ctx context.Context, uid string) error
	Update(ctx context.Context, service *types.Service) (*types.Service, error)

	CreatePort(ctx context.Context, servicePort *types.ServicePort) (string, error)
	GetPorts(ctx context.Context, serviceUid string) ([]*types.ServicePort, error)
	DeletePort(ctx context.Context, uid string) error
	UpdatePort(ctx context.Context, servicePort *types.ServicePort) (*types.ServicePort, error)
}

type KafkaRepository interface {
	Create(ctx context.Context, kafka *types.Kafka) (string, error)
	Get(ctx context.Context, host string) (*types.Kafka, error)
	Delete(ctx context.Context, uid string) error
	Update(ctx context.Context, kafka *types.Kafka) (*types.Kafka, error)

	CreateQueue(ctx context.Context, kafkaQueue *types.KafkaQueue) (string, error)
	GetQueues(ctx context.Context, kafkaUid string) ([]*types.KafkaQueue, error)
	DeleteQueue(ctx context.Context, uid string) error
	UpdateQueue(ctx context.Context, kafkaQueue *types.KafkaQueue) (*types.KafkaQueue, error)
}

type RedisRepository interface {
	Create(ctx context.Context, redis *types.Redis) (string, error)
	Get(ctx context.Context, host string) (*types.Redis, error)
	Delete(ctx context.Context, uid string) error
	Update(ctx context.Context, redis *types.Redis) (*types.Redis, error)

	CreateDB(ctx context.Context, redisDB *types.RedisDB) (string, error)
	GetDBs(ctx context.Context, redisUid string) ([]*types.RedisDB, error)
	DeleteDB(ctx context.Context, uid string) error
	UpdateDB(ctx context.Context, redisDB *types.RedisDB) (*types.RedisDB, error)

	CreateNamespace(ctx context.Context, redisNS *types.RedisNamespace) (string, error)
	GetNamespaces(ctx context.Context, dbUid string) ([]*types.RedisNamespace, error)
	DeleteNamespace(ctx context.Context, uid string) error
	UpdateNamespace(ctx context.Context, redisNS *types.RedisNamespace) (*types.RedisNamespace, error)
}

type RabbitMQRepository interface {
	Create(ctx context.Context, rabbitMQ *types.RabbitMQ) (string, error)
	Get(ctx context.Context, host string) (*types.RabbitMQ, error)
	Delete(ctx context.Context, uid string) error
	Update(ctx context.Context, rabbitMQ *types.RabbitMQ) (*types.RabbitMQ, error)

	CreateQueue(ctx context.Context, rabbitQueue *types.RabbitQueue) (string, error)
	GetQueues(ctx context.Context, rabbitUid string) ([]*types.RabbitQueue, error)
	DeleteQueue(ctx context.Context, uid string) error
	UpdateQueue(ctx context.Context, rabbitQueue *types.RabbitQueue) (*types.RabbitQueue, error)
}

type PostgresRepository interface {
	Create(ctx context.Context, postgres *types.Postgresql) (string, error)
	Get(ctx context.Context, host string) (*types.Postgresql, error)
	Delete(ctx context.Context, uid string) error
	Update(ctx context.Context, postgres *types.Postgresql) (*types.Postgresql, error)

	CreateDB(ctx context.Context, postgresDB *types.PostgresqlDB) (string, error)
	GetDBs(ctx context.Context, postgresUid string) ([]*types.PostgresqlDB, error)
	DeleteDB(ctx context.Context, uid string) error
	UpdateDB(ctx context.Context, postgresDB *types.PostgresqlDB) (*types.PostgresqlDB, error)

	CreateScheme(ctx context.Context, postgresScheme *types.PostgresqlScheme) (string, error)
	GetSchemes(ctx context.Context, dbUid string) ([]*types.PostgresqlScheme, error)
	DeleteScheme(ctx context.Context, uid string) error
	UpdateScheme(ctx context.Context, postgresScheme *types.PostgresqlScheme) (*types.PostgresqlScheme, error)

	CreateTable(ctx context.Context, postgresTable *types.PostgresqlTable) (string, error)
	GetTables(ctx context.Context, schemeUid string) ([]*types.PostgresqlTable, error)
	DeleteTable(ctx context.Context, uid string) error
	UpdateTable(ctx context.Context, postgresTable *types.PostgresqlTable) (*types.PostgresqlTable, error)
}

type ConnectionRepository interface {
	Create(ctx context.Context, fromID, toID *types.ConnNode, connType types.ConnType) error
	Delete(ctx context.Context, fromID, toID *types.ConnNode, connType types.ConnType) error
}

// type TeamRepository interface {
// 	Create(team *types.Team) error
// }

// type PipelineRepository interface {
// 	Create(pipeline *types.Pipeline) error
// }
