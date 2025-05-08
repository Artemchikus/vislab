package neo4j

import (
	"context"
	"fmt"
	"sync"
	"time"
	myConf "vislab/config"
	"vislab/storage"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/log"
)

type Neo4jStorage struct {
	db          neo4j.DriverWithContext
	serviceRepo storage.ServiceRepository
	redisRepo   storage.RedisRepository
	// teamRepo    storage.TeamRepository
	// pipelineRepo storage.PipelineRepository
	connRepo     storage.ConnectionRepository
	postgresRepo storage.PostgresRepository
	kafkaRepo    storage.KafkaRepository
	rabbitRepo   storage.RabbitMQRepository
}

var (
	neo4jInstance *Neo4jStorage
	neo4jOnce     sync.Once
)

func MustStorage(ctx context.Context, myConf *myConf.StorageConfig) (*Neo4jStorage, error) {
	var err error

	neo4jOnce.Do(func() {
		dbUri := fmt.Sprintf("%s://%s:%s", myConf.Name, myConf.Host, myConf.Port)

		logger := func(conf *config.Config) {
			conf.Log = log.ToConsole(log.INFO)
		}

		var db neo4j.DriverWithContext

		db, err = neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(myConf.User, myConf.Password, ""), logger)
		if err != nil {
			return
		}

		for i := 0; i < 10; i++ {
			err = db.VerifyConnectivity(ctx)
			if err == nil {
				break
			}
			time.Sleep(time.Second * 5)
		}
		if err != nil {
			return
		}

		neo4jInstance = &Neo4jStorage{
			db: db,
		}
	})

	return neo4jInstance, err
}

func (n *Neo4jStorage) Disconnect(ctx context.Context) error {
	err := n.db.Close(ctx)
	if err != nil {
		return err
	}

	neo4jOnce = sync.Once{}
	neo4jInstance = nil
	return nil
}
