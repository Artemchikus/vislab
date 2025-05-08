package storefuncs

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"vislab/libs/check"
	"vislab/libs/ptr"
	"vislab/storage"
	storeTypes "vislab/storage/neo4j/types"
	"vislab/types"
)

func storeKafka(ctx context.Context, kafka *types.Kafka, serviceNode *storeTypes.ConnNode, storage storage.Storage) error {
	slog.Info("processing kafka", "kafka", kafka.Host)

	kafkaNode, err := storeKafkaNode(ctx, kafka, storage)
	if err != nil {
		return err
	}

	existingQueues, err := storage.Kafka().GetQueues(ctx, kafkaNode.ID)
	if err != nil {
		return err
	}

	if len(kafka.Queues) == 0 {
		slog.Info("creating dummy kafka queue", "kafka", kafka.Host)
		if err := storeDummyKafkaQueue(ctx, kafkaNode, storage, serviceNode, existingQueues); err != nil {
			return err
		}

		return nil
	}

	for _, queue := range kafka.Queues {
		queueNode, err := storeKafkaQueue(ctx, queue, kafkaNode, existingQueues, storage)
		if err != nil {
			slog.Error("failed to create kafka queue", "queue", queue.Name, "kafka", kafka.Host, "error", err)
			continue
		}

		var connType storeTypes.ConnType
		switch *queue.QueueType {
		case "consumer":
			connType = storeTypes.ConnReceivesFrom
		case "producer":
			connType = storeTypes.ConnSendsTo
		default:
			return fmt.Errorf("unknown queue type: %s", *queue.QueueType)
		}

		slog.Info("creating svc-queue connection", "from_id", serviceNode.ID, "to_id", queueNode.ID, "type", connType)
		if err := storage.Connection().Create(ctx, serviceNode, queueNode, connType); err != nil {
			return err
		}
	}

	return nil
}

func storeKafkaNode(ctx context.Context, kafka *types.Kafka, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeKafka := &storeTypes.Kafka{
		Host: kafka.Host,
		Name: kafka.Name,
		Port: kafka.Port,
	}

	kafkaNode := &storeTypes.ConnNode{
		Class: storeTypes.KafkaClass,
	}

	slog.Info("updating kafka", "kafka", kafka.Host)
	dbKafka, err := storage.Kafka().Update(ctx, storeKafka)
	if err == nil {
		kafkaNode.ID = *dbKafka.UID
		return kafkaNode, nil
	}
	slog.Error("error updating kafka", "kafka", kafka.Host, "error", err)

	if strings.Contains(err.Error(), "nothing to update") {
		dbKafka, err := storage.Kafka().Get(ctx, *storeKafka.Host)
		if err != nil {
			return nil, err
		}

		kafkaNode.ID = *dbKafka.UID
		return kafkaNode, nil
	}

	slog.Info("creating kafka", "kafka", kafka.Host)
	id, err := storage.Kafka().Create(ctx, storeKafka)
	if err != nil {
		return nil, err
	}

	kafkaNode.ID = id
	return kafkaNode, nil
}

func storeKafkaQueue(ctx context.Context, queue *types.KafkaQueue, kafkaNode *storeTypes.ConnNode, existingQueues []*storeTypes.KafkaQueue, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeQueue := &storeTypes.KafkaQueue{
		Name:      queue.Name,
		QueueType: queue.QueueType,
		Topic:     queue.Topic,
		TypeName:  queue.TypeName,
	}

	queueNode := &storeTypes.ConnNode{
		Class: storeTypes.KafkaQueueClass,
	}

	for _, existingQueue := range existingQueues {
		if check.ComparePointers(existingQueue.Name, storeQueue.Name) {
			if existingQueue.Equal(storeQueue) {
				queueNode.ID = *existingQueue.UID
				return queueNode, nil
			}

			storeQueue.UID = existingQueue.UID

			slog.Info("updating kafka queue", "queue", queue.Name)
			dbQueue, err := storage.Kafka().UpdateQueue(ctx, storeQueue)
			if err != nil {
				return nil, err
			}

			queueNode.ID = *dbQueue.UID
			return queueNode, nil
		}
	}

	slog.Info("creating kafka queue", "queue", queue.Name)
	id, err := storage.Kafka().CreateQueue(ctx, storeQueue)
	if err != nil {
		return nil, err
	}

	queueNode.ID = id

	slog.Info("creating kafka-queue connection", "from_id", kafkaNode.ID, "to_id", queueNode.ID, "type", storeTypes.ConnIN)
	if err := storage.Connection().Create(ctx, queueNode, kafkaNode, storeTypes.ConnIN); err != nil {
		return nil, fmt.Errorf("failed to create kafka-queue connection: %w", err)
	}

	return queueNode, nil
}

func storeDummyKafkaQueue(ctx context.Context, kafkaNode *storeTypes.ConnNode, storage storage.Storage, serviceNode *storeTypes.ConnNode, existingQueues []*storeTypes.KafkaQueue) error {
	dummyQueue := &types.KafkaQueue{
		Name: ptr.Ptr("dummy"),
	}

	queueNode, err := storeKafkaQueue(ctx, dummyQueue, kafkaNode, existingQueues, storage)
	if err != nil {
		return err
	}

	slog.Info("creating svc-queue connection", "from_id", serviceNode.ID, "to_id", queueNode.ID, "type", "dummy")
	if err := storage.Connection().Create(ctx, serviceNode, queueNode, "dummy"); err != nil {
		return err
	}

	return nil
}
