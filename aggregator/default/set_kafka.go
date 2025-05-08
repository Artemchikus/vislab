package defaultaggregator

import (
	"context"
	yamlTypes "vislab/sources/yaml/types"
	"vislab/types"
)

func setKafka(ctx context.Context, in []*yamlTypes.Kafka, out *types.All) error {
	for _, kafka := range in {
		newKafka := &types.Kafka{
			Name: kafka.Name,
			Host: kafka.Host,
			Port: kafka.Port,
		}

		for _, queue := range kafka.Queues {
			newQueue := &types.KafkaQueue{
				Name:      queue.Name,
				QueueType: queue.QueueType,
				Topic:     queue.Topic,
				TypeName:  queue.TypeName,
			}

			newKafka.Queues = append(newKafka.Queues, newQueue)
		}

		out.Kafkas = append(out.Kafkas, newKafka)
	}

	return nil
}
