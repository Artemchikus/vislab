package yaml

import (
	"fmt"
	"strconv"
	"strings"
	"vislab/libs/ptr"
	"vislab/sources/yaml/types"
)

func getSetRabbitFunc(pathParts []string) (func(string, *types.All) error, error) {
	switch pathParts[0] {
	case "host":
		return func(s string, all *types.All) error {
			checkRabbit(all)

			host := ptr.Ptr(s)

			if all.RabbitMQ.LastInstance.Host == nil {
				all.RabbitMQ.LastInstance.Host = host
				return nil
			}

			rabbit := &types.RabbitMQ{Host: host}
			all.RabbitMQ.Instances = append(all.RabbitMQ.Instances, rabbit)
			all.RabbitMQ.LastInstance = rabbit

			return nil
		}, nil
	case "port":
		return func(s string, all *types.All) error {
			checkRabbit(all)

			intVal, err := strconv.Atoi(s)
			if err != nil {
				return err
			}

			port := ptr.Ptr(int64(intVal))

			if all.RabbitMQ.LastInstance.Port == nil {
				all.RabbitMQ.LastInstance.Port = port
				return nil
			}

			rabbit := &types.RabbitMQ{Port: port}
			all.RabbitMQ.Instances = append(all.RabbitMQ.Instances, rabbit)
			all.RabbitMQ.LastInstance = rabbit

			return nil
		}, nil
	case "user":
		return func(s string, all *types.All) error {
			checkRabbit(all)

			user := ptr.Ptr(s)

			if all.RabbitMQ.LastInstance.User == nil {
				all.RabbitMQ.LastInstance.User = user
				return nil
			}

			rabbit := &types.RabbitMQ{User: user}
			all.RabbitMQ.Instances = append(all.RabbitMQ.Instances, rabbit)
			all.RabbitMQ.LastInstance = rabbit

			return nil
		}, nil
	case "queue":
		if len(pathParts) < 2 {
			return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
		}

		switch pathParts[1] {
		case "name":
			return func(s string, all *types.All) error {
				checkRabbitQueues(all)

				name := ptr.Ptr(s)

				if all.RabbitMQ.LastInstance.LastQueue.Name == nil {
					all.RabbitMQ.LastInstance.LastQueue.Name = name
					return nil
				}

				queue := &types.RabbitQueue{Name: name}
				all.RabbitMQ.LastInstance.Queues = append(all.RabbitMQ.LastInstance.Queues, queue)
				all.RabbitMQ.LastInstance.LastQueue = queue

				return nil
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
}

func checkRabbit(all *types.All) {
	if all.RabbitMQ == nil {
		rabbit := &types.RabbitMQ{}
		all.RabbitMQ = &types.RabbitMQs{
			Instances:    []*types.RabbitMQ{rabbit},
			LastInstance: rabbit,
		}
	}
}

func checkRabbitQueues(all *types.All) {
	checkRabbit(all)

	if all.RabbitMQ.LastInstance.Queues == nil {
		queue := &types.RabbitQueue{}
		all.RabbitMQ.LastInstance.LastQueue = queue
		all.RabbitMQ.LastInstance.Queues = []*types.RabbitQueue{queue}
	}
}
