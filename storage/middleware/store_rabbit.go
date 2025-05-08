package storefuncs

import (
	"context"
	"log/slog"
	"strings"
	"vislab/libs/check"
	"vislab/libs/ptr"
	"vislab/storage"
	storeTypes "vislab/storage/neo4j/types"
	"vislab/types"
)

func storeRabbitMQ(ctx context.Context, rabbitMQ *types.RabbitMQ, serviceNode *storeTypes.ConnNode, storage storage.Storage) error {
	rabbitMQNode, err := storeRabbitMQNode(ctx, rabbitMQ, serviceNode, storage)
	if err != nil {
		return err
	}

	slog.Info("getting rabbitmq queues", "rabbitmq", rabbitMQ.Host)
	existingQueues, err := storage.RabbitMQ().GetQueues(ctx, rabbitMQNode.ID)
	if err != nil {
		return err
	}

	if len(rabbitMQ.Queues) == 0 {
		slog.Info("creating dummy rabbitmq queue", "rabbitmq", rabbitMQ.Host)
		if err := storeDummyRabbitMQQueue(ctx, rabbitMQNode, storage, serviceNode, existingQueues); err != nil {
			return err
		}

		return nil
	}

	for _, queue := range rabbitMQ.Queues {
		queueNode, err := storeRabbitMQQueue(ctx, queue, rabbitMQNode, existingQueues, storage)
		if err != nil {
			slog.Error("failed to create rabbitmq queue", "error", err)
			continue
		}

		slog.Info("creating svc-queue connection", "from_id", serviceNode.ID, "to_id", queueNode.ID, "type", storeTypes.ConnUses)
		if err := storage.Connection().Create(ctx, serviceNode, queueNode, storeTypes.ConnUses); err != nil {
			return err
		}
	}

	return nil
}

func storeRabbitMQNode(ctx context.Context, rabbitMQ *types.RabbitMQ, serviceNode *storeTypes.ConnNode, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeRabbitMQ := &storeTypes.RabbitMQ{
		Host: rabbitMQ.Host,
		Port: rabbitMQ.Port,
		User: rabbitMQ.User,
	}

	rabbitMQNode := &storeTypes.ConnNode{
		Class: storeTypes.RabbitMQClass,
	}

	slog.Info("updating rabbitmq", "rabbitmq", rabbitMQ.Host)
	dbRabbitMQ, err := storage.RabbitMQ().Update(ctx, storeRabbitMQ)
	if err == nil {
		rabbitMQNode.ID = *dbRabbitMQ.UID
		return rabbitMQNode, nil
	}
	slog.Error("error updating rabbitmq", "rabbitmq", rabbitMQ.Host, "error", err)

	if strings.Contains(err.Error(), "nothing to update") {
		dbRabbitMQ, err := storage.RabbitMQ().Get(ctx, *storeRabbitMQ.Host)
		if err != nil {
			return nil, err
		}

		rabbitMQNode.ID = *dbRabbitMQ.UID
		return rabbitMQNode, nil
	}

	slog.Info("creating rabbitmq", "rabbitmq", rabbitMQ.Host)
	id, err := storage.RabbitMQ().Create(ctx, storeRabbitMQ)
	if err != nil {
		return nil, err
	}

	rabbitMQNode.ID = id
	return rabbitMQNode, nil
}

func storeRabbitMQQueue(ctx context.Context, queue *types.RabbitQueue, rabbitMQNode *storeTypes.ConnNode, existingQueues []*storeTypes.RabbitQueue, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeQueue := &storeTypes.RabbitQueue{
		Name: queue.Name,
	}

	queueNode := &storeTypes.ConnNode{
		Class: storeTypes.RabbitQueueClass,
	}

	for _, existingQueue := range existingQueues {
		if check.ComparePointers(existingQueue.Name, storeQueue.Name) {
			if existingQueue.Equal(storeQueue) {
				queueNode.ID = *existingQueue.UID
				return queueNode, nil
			}

			storeQueue.UID = existingQueue.UID

			slog.Info("updating rabbitmq queue", "queue", queue.Name)
			dbQueue, err := storage.RabbitMQ().UpdateQueue(ctx, storeQueue)
			if err != nil {
				return nil, err
			}

			queueNode.ID = *dbQueue.UID
			return queueNode, nil
		}
	}

	slog.Info("creating rabbitmq queue", "queue", queue.Name)
	id, err := storage.RabbitMQ().CreateQueue(ctx, storeQueue)
	if err != nil {
		return nil, err
	}

	queueNode.ID = id

	slog.Info("creating rabbitmq-queue connection", "from_id", queueNode.ID, "to_id", rabbitMQNode.ID, "type", storeTypes.ConnIN)
	if err := storage.Connection().Create(ctx, queueNode, rabbitMQNode, storeTypes.ConnIN); err != nil {
		return nil, err
	}

	return queueNode, nil
}

func storeDummyRabbitMQQueue(ctx context.Context, rabbitMQNode *storeTypes.ConnNode, storage storage.Storage, serviceNode *storeTypes.ConnNode, existingQueues []*storeTypes.RabbitQueue) error {
	dummyQueue := &types.RabbitQueue{
		Name: ptr.Ptr("dummy"),
	}

	queueNode, err := storeRabbitMQQueue(ctx, dummyQueue, rabbitMQNode, existingQueues, storage)
	if err != nil {
		return err
	}

	slog.Info("creating svc-queue connection", "from_id", serviceNode.ID, "to_id", queueNode.ID, "type", "dummy")
	if err := storage.Connection().Create(ctx, serviceNode, queueNode, "dummy"); err != nil {
		return err
	}

	return nil
}
