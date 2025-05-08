package updater

import "context"

type (
	Updater interface {
		Start(ctx context.Context) error
	}
)
