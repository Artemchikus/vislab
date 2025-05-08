package storefuncs

import (
	"context"
	"log/slog"
	"strings"
	"vislab/libs/check"
	"vislab/storage"
	storeTypes "vislab/storage/neo4j/types"
	"vislab/types"
)

func storeService(ctx context.Context, service *types.Service, storage storage.Storage) (*storeTypes.ConnNode, error) {
	serviceNode, err := storeServiceNode(ctx, service, storage)
	if err != nil {
		return nil, err
	}

	existingPorts, err := storage.Service().GetPorts(ctx, serviceNode.ID)
	if err != nil {
		return nil, err
	}

	for _, port := range service.Ports {
		_, err := storeServicePort(ctx, port, serviceNode, existingPorts, storage)
		if err != nil {
			return nil, err
		}
	}

	return serviceNode, nil
}

func storeServiceNode(ctx context.Context, service *types.Service, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storeService := &storeTypes.Service{
		Name:        service.Name,
		FullName:    service.FullName,
		Link:        service.Link,
		MainBranch:  service.MainBranch,
		LatestTag:   service.LatestTag,
		Language:    service.Language,
		Description: service.Description,
		Status:      service.Status,
	}

	serviceNode := &storeTypes.ConnNode{
		Class: storeTypes.ServiceClass,
	}

	slog.Info("updating service", "service", service.Name)
	dbService, err := storage.Service().Update(ctx, storeService)
	if err == nil {
		serviceNode.ID = *dbService.UID
		return serviceNode, nil
	}
	slog.Warn("error updating service", "service", *service.Name, "error", err)

	if strings.Contains(err.Error(), "nothing to update") {
		dbService, err := storage.Service().Get(ctx, *storeService.Name)
		if err != nil {
			return nil, err
		}

		serviceNode.ID = *dbService.UID
		return serviceNode, nil
	}

	slog.Info("creating service", "service", service.Name)
	id, err := storage.Service().Create(ctx, storeService)
	if err != nil {
		return nil, err
	}

	serviceNode.ID = id
	return serviceNode, nil
}

func storeServicePort(ctx context.Context, port *types.Port, serviceNode *storeTypes.ConnNode, existingPorts []*storeTypes.ServicePort, storage storage.Storage) (*storeTypes.ConnNode, error) {
	storePort := &storeTypes.ServicePort{
		Number: port.Number,
	}

	portNode := &storeTypes.ConnNode{
		Class: storeTypes.ServicePortClass,
	}

	for _, existingPort := range existingPorts {
		if check.ComparePointers(existingPort.Number, storePort.Number) {
			if existingPort.Equal(storePort) {
				portNode.ID = *existingPort.UID
				return portNode, nil
			}

			storePort.UID = existingPort.UID

			slog.Info("updating port", "port", *storePort.Number)
			dbPort, err := storage.Service().UpdatePort(ctx, storePort)
			if err != nil {
				return nil, err
			}

			portNode.ID = *dbPort.UID
			return portNode, nil
		}
	}

	slog.Info("creating port", "port", *storePort.Number)
	id, err := storage.Service().CreatePort(ctx, storePort)
	if err != nil {
		return nil, err
	}

	portNode.ID = id

	slog.Info("creating svc-port connection", "from_id", portNode.ID, "to_id", serviceNode.ID, "type", storeTypes.ConnIN)
	if err := storage.Connection().Create(ctx, portNode, serviceNode, storeTypes.ConnIN); err != nil {
		return nil, err
	}

	return portNode, nil
}

func storeOtherService(ctx context.Context, otherService *types.Service, serviceNode *storeTypes.ConnNode, storage storage.Storage) error {
	otherServiceNode, err := storeServiceNode(ctx, otherService, storage)
	if err != nil {
		return err
	}

	existingPorts, err := storage.Service().GetPorts(ctx, otherServiceNode.ID)
	if err != nil {
		return err
	}

	for _, port := range otherService.Ports {
		portNode, err := storeServicePort(ctx, port, otherServiceNode, existingPorts, storage)
		if err != nil {
			slog.Error("error storing other-svc-port", "error", err)
			continue
		}

		slog.Info("creating svc-other-svc-port connection", "from_id", serviceNode.ID, "to_id", portNode.ID, "type", storeTypes.ConnUses)

		if err := storage.Connection().Create(ctx, serviceNode, portNode, storeTypes.ConnUses); err != nil {
			return err
		}
	}

	return nil
}
