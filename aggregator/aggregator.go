package aggregator

import (
	"context"
	"vislab/types"
)

type (
	Aggregator interface {
		Set(context.Context, any) error
		Get(context.Context) (*types.All, error)
	}
)
