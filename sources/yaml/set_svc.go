package yaml

import (
	"fmt"
	"strconv"
	"strings"
	"vislab/libs/ptr"
	"vislab/sources/yaml/types"
)

func getSetSvcFunc(pathParts []string) (func(string, *types.All) error, error) {
	switch pathParts[0] {
	case "name":
		return func(s string, all *types.All) error {
			checkSvc(all)

			name := ptr.Ptr(s)

			if all.Service.LastInstance.Name == nil {
				all.Service.LastInstance.Name = name
				return nil
			}

			service := &types.Service{Name: name}
			all.Service.Instances = append(all.Service.Instances, service)
			all.Service.LastInstance = service

			return nil
		}, nil
	case "full_name":
		return func(s string, all *types.All) error {
			checkSvc(all)

			fullName := ptr.Ptr(s)

			if all.Service.LastInstance.FullName == nil {
				all.Service.LastInstance.FullName = fullName
				return nil
			}

			service := &types.Service{FullName: fullName}
			all.Service.Instances = append(all.Service.Instances, service)
			all.Service.LastInstance = service

			return nil
		}, nil
	case "project_id":
		return func(s string, all *types.All) error {
			checkSvc(all)

			intVal, err := strconv.Atoi(s)
			if err != nil {
				return err
			}

			ProjectID := ptr.Ptr(int64(intVal))

			if all.Service.LastInstance.ProjectID == nil {
				all.Service.LastInstance.ProjectID = ProjectID
				return nil
			}

			service := &types.Service{ProjectID: ProjectID}
			all.Service.Instances = append(all.Service.Instances, service)
			all.Service.LastInstance = service

			return nil
		}, nil
	case "tag":
		return func(s string, all *types.All) error {
			checkSvc(all)

			tag := ptr.Ptr(s)

			if all.Service.LastInstance.Tag == nil {
				all.Service.LastInstance.Tag = tag
				return nil
			}

			service := &types.Service{Tag: tag}
			all.Service.Instances = append(all.Service.Instances, service)
			all.Service.LastInstance = service

			return nil
		}, nil
	case "port":
		if len(pathParts) < 2 {
			return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
		}

		switch pathParts[1] {
		case "number":
			return func(s string, all *types.All) error {
				checkSvcPorts(all)

				intVal, err := strconv.Atoi(s)
				if err != nil {
					return err
				}

				number := ptr.Ptr(int64(intVal))

				if all.Service.LastInstance.LastPort.Number == nil {
					all.Service.LastInstance.LastPort.Number = number
					return nil
				}

				port := &types.Port{Number: number}
				all.Service.LastInstance.Ports = append(all.Service.LastInstance.Ports, port)
				all.Service.LastInstance.LastPort = port

				return nil
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
}

func checkSvc(all *types.All) {
	if all.Service == nil {
		svc := &types.Service{}
		all.Service = &types.Services{
			Instances:    []*types.Service{svc},
			LastInstance: svc,
		}
	}
}

func checkSvcPorts(all *types.All) {
	checkSvc(all)

	if all.Service.LastInstance.Ports == nil {
		port := &types.Port{}
		all.Service.LastInstance.LastPort = port
		all.Service.LastInstance.Ports = []*types.Port{port}
	}
}
