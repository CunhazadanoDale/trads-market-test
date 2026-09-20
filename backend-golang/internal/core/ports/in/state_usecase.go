package in

import "context"

type StateUseCase interface {
	Import(ctx context.Context) error
}