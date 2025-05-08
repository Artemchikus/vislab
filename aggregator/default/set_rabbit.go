package defaultaggregator

import (
	"context"
	yamlTypes "vislab/sources/yaml/types"
	"vislab/types"
)

func setRabbit(ctx context.Context, in []*yamlTypes.RabbitMQ, out *types.All) error {
	for _, rabbitMQ := range in {
		newRabbitMQ := &types.RabbitMQ{
			Host: rabbitMQ.Host,
			Port: rabbitMQ.Port,
			User: rabbitMQ.User,
		}

		for _, queue := range rabbitMQ.Queues {
			newQueue := &types.RabbitQueue{
				Name: queue.Name,
				// QueueType: queue.QueueType,
				// Topic:     queue.Topic,
				// TypeName:  queue.TypeName,
			}

			newRabbitMQ.Queues = append(newRabbitMQ.Queues, newQueue)
		}

		out.RabbitMQs = append(out.RabbitMQs, newRabbitMQ)
	}

	return nil
}
