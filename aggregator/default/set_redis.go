package defaultaggregator

import (
	"context"
	yamlTypes "vislab/sources/yaml/types"
	"vislab/types"
)

func setRedis(ctx context.Context, in []*yamlTypes.Redis, out *types.All) error {
	for _, redis := range in {
		newRedis := &types.Redis{
			Host:   redis.Host,
			Port:   redis.Port,
			Master: redis.Master,
		}

		for _, database := range redis.Databases {
			newDatabase := &types.RedisDB{
				Name: database.Name,
				// Owner: database.Owner,
			}

			for _, namespaces := range database.Namespaces {
				newDatabase.Namespaces = append(newDatabase.Namespaces, &types.RedisNamespace{
					Name: namespaces.Name,
				})
			}

			newRedis.Databases = append(newRedis.Databases, newDatabase)
		}

		out.Redises = append(out.Redises, newRedis)
	}

	return nil
}
