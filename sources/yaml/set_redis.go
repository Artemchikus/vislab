package yaml

import (
	"fmt"
	"strconv"
	"strings"
	"vislab/libs/ptr"
	"vislab/sources/yaml/types"
)

func getSetRedisFunc(pathParts []string) (func(string, *types.All) error, error) {
	switch pathParts[0] {
	case "host":
		return func(s string, all *types.All) error {
			checkRedis(all)

			host := ptr.Ptr(s)

			if all.Redis.LastInstance.Host == nil {
				all.Redis.LastInstance.Host = host
				return nil
			}

			redis := &types.Redis{Host: host}
			all.Redis.Instances = append(all.Redis.Instances, redis)
			all.Redis.LastInstance = redis

			return nil
		}, nil
	case "port":
		return func(s string, all *types.All) error {
			checkRedis(all)

			intVal, err := strconv.Atoi(s)
			if err != nil {
				return err
			}

			port := ptr.Ptr(int64(intVal))

			if all.Redis.LastInstance.Port == nil {
				all.Redis.LastInstance.Port = port
				return nil
			}

			redis := &types.Redis{Port: port}
			all.Redis.Instances = append(all.Redis.Instances, redis)
			all.Redis.LastInstance = redis

			return nil
		}, nil
	case "master":
		return func(s string, all *types.All) error {
			checkRedis(all)

			master := ptr.Ptr(s)

			if all.Redis.LastInstance.Master == nil {
				all.Redis.LastInstance.Master = master
				return nil
			}

			redis := &types.Redis{Master: master}
			all.Redis.Instances = append(all.Redis.Instances, redis)
			all.Redis.LastInstance = redis

			return nil
		}, nil
	case "sentinel":
		if len(pathParts) < 2 {
			return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
		}

		switch pathParts[1] {
		case "host":
			return func(s string, all *types.All) error {
				checkRedisSentinel(all)

				host := ptr.Ptr(s)

				if all.Redis.LastInstance.Sentinel.Host == nil {
					all.Redis.LastInstance.Sentinel.Host = host
					return nil
				}

				sentinel := &types.Sentinel{Host: host}
				all.Redis.LastInstance.Sentinel = sentinel

				return nil
			}, nil
		case "port":
			return func(s string, all *types.All) error {
				checkRedisSentinel(all)

				intVal, err := strconv.Atoi(s)
				if err != nil {
					return err
				}

				port := ptr.Ptr(int64(intVal))

				if all.Redis.LastInstance.Sentinel.Port == nil {
					all.Redis.LastInstance.Sentinel.Port = port
					return nil
				}

				sentiel := &types.Sentinel{Port: port}
				all.Redis.LastInstance.Sentinel = sentiel

				return nil
			}, nil
		}
	case "database":
		if len(pathParts) < 2 {
			return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
		}

		switch pathParts[1] {
		case "name":
			return func(s string, all *types.All) error {
				checkRedisDb(all)

				name := ptr.Ptr(s)

				if all.Redis.LastInstance.LastDatabase.Name == nil {
					all.Redis.LastInstance.LastDatabase.Name = name
					return nil
				}

				database := &types.RedisDB{Name: name}
				all.Redis.LastInstance.Databases = append(all.Redis.LastInstance.Databases, database)
				all.Redis.LastInstance.LastDatabase = database

				return nil
			}, nil
		case "namespace":
			if len(pathParts) < 3 {
				return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
			}

			switch pathParts[2] {
			case "name":
				return func(s string, all *types.All) error {
					checkRedisNamespace(all)

					name := ptr.Ptr(s)

					if all.Redis.LastInstance.LastDatabase.LastNamespace.Name == nil {
						all.Redis.LastInstance.LastDatabase.LastNamespace.Name = name
						return nil
					}

					namespace := &types.RedisNamespace{Name: name}
					all.Redis.LastInstance.LastDatabase.Namespaces = append(all.Redis.LastInstance.LastDatabase.Namespaces, namespace)
					all.Redis.LastInstance.LastDatabase.LastNamespace = namespace

					return nil
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
}

func checkRedis(all *types.All) {
	if all.Redis == nil {
		redis := &types.Redis{}
		all.Redis = &types.Redises{
			Instances:    []*types.Redis{redis},
			LastInstance: redis,
		}
	}
}

func checkRedisSentinel(all *types.All) {
	checkRedis(all)

	if all.Redis.LastInstance.Sentinel == nil {
		all.Redis.LastInstance.Sentinel = &types.Sentinel{}
	}
}

func checkRedisDb(all *types.All) {
	checkRedis(all)

	if all.Redis.LastInstance.Databases == nil {
		db := &types.RedisDB{}
		all.Redis.LastInstance.LastDatabase = db
		all.Redis.LastInstance.Databases = []*types.RedisDB{db}
	}
}

func checkRedisNamespace(all *types.All) {
	checkRedisDb(all)

	if all.Redis.LastInstance.LastDatabase.Namespaces == nil {
		ns := &types.RedisNamespace{}
		all.Redis.LastInstance.LastDatabase.LastNamespace = ns
		all.Redis.LastInstance.LastDatabase.Namespaces = []*types.RedisNamespace{ns}
	}
}
