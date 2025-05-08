package defaultaggregator

import (
	"context"
	yamlTypes "vislab/sources/yaml/types"
	"vislab/types"
)

func setOtherService(ctx context.Context, in []*yamlTypes.Service, out *types.All) error {
	for _, service := range in {
		newOtherService := &types.Service{
			Name: service.Name,
			// Description: service.Description,
			// Group:      service.Group,
			// Language:    service.Language,
			// Link:        service.Link,
			// MainBranch:  service.MainBranch,
		}

		for _, port := range service.Ports {
			newOtherService.Ports = append(newOtherService.Ports, &types.Port{
				Number: port.Number,
			})
		}

		out.OtherServices = append(out.OtherServices, newOtherService)
	}
	return nil
}
