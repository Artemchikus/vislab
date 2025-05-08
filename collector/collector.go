package collector

import (
	"context"
)

type CollectorOption func(Collector) error

type Collector interface {
	Collect(ctx context.Context) error
	Update(ctx context.Context, options ...CollectorOption) error
}
