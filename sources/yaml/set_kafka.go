package yaml

import (
	"fmt"
	"strconv"
	"strings"
	"vislab/libs/ptr"
	"vislab/sources/yaml/types"
)

func getSetKafkaFunc(pathParts []string) (func(string, *types.All) error, error) {
	switch pathParts[0] {
	case "name":
		return func(s string, all *types.All) error {
			checkKafka(all)

			name := ptr.Ptr(s)

			if all.Kafka.LastInstance.Name == nil {
				all.Kafka.LastInstance.Name = name
				return nil
			}

			kafka := &types.Kafka{Name: name}
			all.Kafka.Instances = append(all.Kafka.Instances, kafka)
			all.Kafka.LastInstance = kafka

			return nil
		}, nil
	case "host":
		return func(s string, all *types.All) error {
			checkKafka(all)

			host := ptr.Ptr(s)

			if all.Kafka.LastInstance.Host == nil {
				all.Kafka.LastInstance.Host = host
				return nil
			}

			kafka := &types.Kafka{Host: host}
			all.Kafka.Instances = append(all.Kafka.Instances, kafka)
			all.Kafka.LastInstance = kafka

			return nil
		}, nil
	case "port":
		return func(s string, all *types.All) error {
			checkKafka(all)

			intVal, err := strconv.Atoi(s)
			if err != nil {
				return err
			}

			port := ptr.Ptr(int64(intVal))

			if all.Kafka.LastInstance.Port == nil {
				all.Kafka.LastInstance.Port = port
				return nil
			}

			kafka := &types.Kafka{Port: port}
			all.Kafka.Instances = append(all.Kafka.Instances, kafka)
			all.Kafka.LastInstance = kafka

			return nil
		}, nil
	case "queue":
		if len(pathParts) < 2 {
			return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
		}

		switch pathParts[1] {
		case "name":
			return func(s string, all *types.All) error {
				checkKafkaQueues(all)

				name := ptr.Ptr(s)

				if all.Kafka.LastInstance.LastQueue.Name == nil {
					all.Kafka.LastInstance.LastQueue.Name = name
					return nil
				}

				queue := &types.KafkaQueue{Name: name}
				all.Kafka.LastInstance.Queues = append(all.Kafka.LastInstance.Queues, queue)
				all.Kafka.LastInstance.LastQueue = queue

				return nil
			}, nil
		case "type":
			return func(s string, all *types.All) error {
				checkKafkaQueues(all)

				ty := ptr.Ptr(s)

				if all.Kafka.LastInstance.LastQueue.QueueType == nil {
					all.Kafka.LastInstance.LastQueue.QueueType = ty
					return nil
				}

				queue := &types.KafkaQueue{QueueType: ty}
				all.Kafka.LastInstance.Queues = append(all.Kafka.LastInstance.Queues, queue)
				all.Kafka.LastInstance.LastQueue = queue

				return nil
			}, nil
		case "topic":
			return func(s string, all *types.All) error {
				checkKafkaQueues(all)

				topic := ptr.Ptr(s)

				if all.Kafka.LastInstance.LastQueue.Topic == nil {
					all.Kafka.LastInstance.LastQueue.Topic = topic
					return nil
				}

				queue := &types.KafkaQueue{Topic: topic}
				all.Kafka.LastInstance.Queues = append(all.Kafka.LastInstance.Queues, queue)
				all.Kafka.LastInstance.LastQueue = queue

				return nil
			}, nil
		case "type_name":
			return func(s string, all *types.All) error {
				checkKafkaQueues(all)

				typeName := ptr.Ptr(s)

				if all.Kafka.LastInstance.LastQueue.TypeName == nil {
					all.Kafka.LastInstance.LastQueue.TypeName = typeName
					return nil
				}

				queue := &types.KafkaQueue{TypeName: typeName}
				all.Kafka.LastInstance.Queues = append(all.Kafka.LastInstance.Queues, queue)
				all.Kafka.LastInstance.LastQueue = queue

				return nil
			}, nil
		}
	}
	return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
}

func checkKafka(all *types.All) {
	if all.Kafka == nil {
		kafka := &types.Kafka{}
		all.Kafka = &types.Kafkas{
			Instances:    []*types.Kafka{kafka},
			LastInstance: kafka,
		}
	}
}

func checkKafkaQueues(all *types.All) {
	checkKafka(all)

	if all.Kafka.LastInstance.Queues == nil {
		queue := &types.KafkaQueue{}
		all.Kafka.LastInstance.LastQueue = queue
		all.Kafka.LastInstance.Queues = []*types.KafkaQueue{queue}
	}
}
