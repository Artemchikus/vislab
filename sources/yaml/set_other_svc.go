package yaml

import (
	"fmt"
	"strconv"
	"strings"
	"vislab/libs/ptr"
	"vislab/sources/yaml/types"
)

func getSetOtherSvcFunc(pathParts []string) (func(string, *types.All) error, error) {
	switch pathParts[0] {
	case "name":
		return func(s string, all *types.All) error {
			checkOtherSvc(all)

			name := ptr.Ptr(s)

			if all.OtherService.LastInstance.Name == nil {
				all.OtherService.LastInstance.Name = name
				return nil
			}

			service := &types.Service{Name: name}
			all.OtherService.Instances = append(all.OtherService.Instances, service)
			all.OtherService.LastInstance = service

			return nil
		}, nil
	case "port":
		if len(pathParts) < 2 {
			return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
		}

		switch pathParts[1] {
		case "number":
			return func(s string, all *types.All) error {
				checkOtherSvcPorts(all)

				intVal, err := strconv.Atoi(s)
				if err != nil {
					return err
				}

				number := ptr.Ptr(int64(intVal))

				if all.OtherService.LastInstance.LastPort.Number == nil {
					all.OtherService.LastInstance.LastPort.Number = number
					return nil
				}

				port := &types.Port{Number: number}
				all.OtherService.LastInstance.Ports = append(all.OtherService.LastInstance.Ports, port)
				all.OtherService.LastInstance.LastPort = port

				return nil
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid obj path %s", strings.Join(pathParts, "."))
}

func checkOtherSvc(all *types.All) {
	if all.OtherService == nil {
		srv := &types.Service{}
		all.OtherService = &types.OtherServices{
			Instances:    []*types.Service{srv},
			LastInstance: srv,
		}
	}
}

func checkOtherSvcPorts(all *types.All) {
	checkOtherSvc(all)

	if all.OtherService.LastInstance.Ports == nil {
		port := &types.Port{}
		all.OtherService.LastInstance.LastPort = port
		all.OtherService.LastInstance.Ports = []*types.Port{port}
	}
}
