package in

import "context"

type GDPUsecase interface {
	Import(ctx context.Context) error
}
