package yaml

import (
	"fmt"
	"strconv"
	"strings"
	"vislab/libs/ptr"
	"vislab/sources/yaml/types"
)

func getSetPqFunc(pathParts []string) (func(string, *types.All) error, error) {
	switch pathParts[0] {
	case "host":
		return func(s string, all *types.All) error {
			checkPq(all)

			host := ptr.Ptr(s)

			if all.Postgresql.LastInstance.Host == nil {
				all.Postgresql.LastInstance.Host = host
				return nil
			}

			psql := &types.Postgresql{Host: host}
			all.Postgresql.Instances = append(all.Postgresql.Instances, psql)
			all.Postgresql.LastInstance = psql

			return nil
		}, nil
	case "port":
		return func(s string, all *types.All) error {
			checkPq(all)

			intVal, err := strconv.Atoi(s)
			if err != nil {
				return err
			}

			port := ptr.Ptr(int64(intVal))

			if all.Postgresql.LastInstance.Port == nil {
				all.Postgresql.LastInstance.Port = port
				return nil
			}

			psql := &types.Postgresql{Port: port}
			all.Postgresql.Instances = append(all.Postgresql.Instances, psql)
			all.Postgresql.LastInstance = psql

			return nil
		}, nil
	case "database":
		if len(pathParts) < 2 {
			return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
		}

		switch pathParts[1] {
		case "name":
			return func(s string, all *types.All) error {
				checkPqDBs(all)

				name := ptr.Ptr(s)

				if all.Postgresql.LastInstance.LastDatabase.Name == nil {
					all.Postgresql.LastInstance.LastDatabase.Name = name
					return nil
				}

				db := &types.PqDB{Name: name}
				all.Postgresql.LastInstance.Databases = append(all.Postgresql.LastInstance.Databases, db)
				all.Postgresql.LastInstance.LastDatabase = db

				return nil
			}, nil
		case "for_migrations":
			return func(s string, all *types.All) error {
				checkPqDBs(all)

				boolVal, err := strconv.ParseBool(s)
				if err != nil {
					return err
				}

				forMigrations := ptr.Ptr(boolVal)

				if all.Postgresql.LastInstance.LastDatabase.ForMigrations == nil {
					all.Postgresql.LastInstance.LastDatabase.ForMigrations = forMigrations
					return nil
				}

				db := &types.PqDB{ForMigrations: forMigrations}
				all.Postgresql.LastInstance.Databases = append(all.Postgresql.LastInstance.Databases, db)
				all.Postgresql.LastInstance.LastDatabase = db

				return nil
			}, nil
		case "scheme":
			if len(pathParts) < 3 {
				return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
			}

			switch pathParts[2] {
			case "name":
				return func(s string, all *types.All) error {
					checkPqSchemes(all)

					name := ptr.Ptr(s)

					if all.Postgresql.LastInstance.LastDatabase.LastScheme.Name == nil {
						all.Postgresql.LastInstance.LastDatabase.LastScheme.Name = name
						return nil
					}

					scheme := &types.PqScheme{Name: name}
					all.Postgresql.LastInstance.LastDatabase.Schemes = append(all.Postgresql.LastInstance.LastDatabase.Schemes, scheme)
					all.Postgresql.LastInstance.LastDatabase.LastScheme = scheme

					return nil
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
}

func checkPq(all *types.All) {
	if all.Postgresql == nil {
		psql := &types.Postgresql{}
		all.Postgresql = &types.Postgresqls{
			Instances:    []*types.Postgresql{psql},
			LastInstance: psql,
		}
	}
}

func checkPqDBs(all *types.All) {
	checkPq(all)

	if all.Postgresql.LastInstance.Databases == nil {
		db := &types.PqDB{}
		all.Postgresql.LastInstance.LastDatabase = db
		all.Postgresql.LastInstance.Databases = []*types.PqDB{db}
	}
}

func checkPqSchemes(all *types.All) {
	checkPqDBs(all)

	if all.Postgresql.LastInstance.LastDatabase.Schemes == nil {
		scheme := &types.PqScheme{}
		all.Postgresql.LastInstance.LastDatabase.LastScheme = scheme
		all.Postgresql.LastInstance.LastDatabase.Schemes = []*types.PqScheme{scheme}
	}
}
