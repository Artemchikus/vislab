package gtlabjobsteps

import (
	"context"
	"vislab/aggregator"
)

type (
	Step interface {
		Run(ctx context.Context, params *StepParams) error
		Weight() int64
	}

	StepParams struct {
		ServiceId  int64
		ServiceRef string
		Aggregator aggregator.Aggregator
	}
)
