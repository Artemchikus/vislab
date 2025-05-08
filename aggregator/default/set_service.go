package defaultaggregator

import (
	"context"
	"fmt"
	yamlTypes "vislab/sources/yaml/types"
	"vislab/types"
)

func setService(ctx context.Context, in []*yamlTypes.Service, out *types.All) error {
	if len(in) == 0 {
		return fmt.Errorf("no service data provided")
	}

	if in[0].Name != nil {
		out.Service.Name = in[0].Name
	}

	for _, port := range in[0].Ports {
		out.Service.Ports = append(out.Service.Ports, &types.Port{
			Number: port.Number,
		})
	}
	return nil
}
