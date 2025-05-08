package storefuncs

import (
	"context"
	"log/slog"
	"vislab/storage"
	"vislab/types"
)

func StoreResources(ctx context.Context, resInfo *types.All, storage storage.Storage) error {
	serviceNode, err := storeService(ctx, resInfo.Service, storage)
	if err != nil {
		return err
	}

	if len(resInfo.Kafkas) == 0 {
		slog.Debug("no kafkas found in resource yaml")
	} else {
		for _, kafka := range resInfo.Kafkas {
			if err := storeKafka(ctx, kafka, serviceNode, storage); err != nil {
				slog.Error("failed to store kafkas", "err", err)
				continue
			}
		}
	}

	if len(resInfo.RabbitMQs) == 0 {
		slog.Debug("no rabbitmq found in resource yaml")
	} else {
		for _, rabbitmq := range resInfo.RabbitMQs {
			if err := storeRabbitMQ(ctx, rabbitmq, serviceNode, storage); err != nil {
				slog.Error("failed to store rabbitmq", "err", err)
				continue
			}
		}
	}

	if len(resInfo.Postgresqls) == 0 {
		slog.Debug("no postgresql found in resource yaml")
	} else {
		for _, postgresql := range resInfo.Postgresqls {
			if err := storePostgres(ctx, postgresql, serviceNode, storage); err != nil {
				slog.Error("failed to store postgresql", "err", err)
				continue
			}
		}
	}

	if len(resInfo.Redises) == 0 {
		slog.Debug("no redis found in resource yaml")
	} else {
		for _, redis := range resInfo.Redises {
			if err := storeRedis(ctx, redis, serviceNode, storage); err != nil {
				slog.Error("failed to store redis", "err", err)
				continue
			}
		}
	}

	if len(resInfo.OtherServices) == 0 {
		slog.Debug("no other services found in resource yaml")
	} else {
		for _, otherService := range resInfo.OtherServices {
			if err := storeOtherService(ctx, otherService, serviceNode, storage); err != nil {
				slog.Error("failed to store other service", "err", err)
				continue
			}
		}
	}

	return nil
}
