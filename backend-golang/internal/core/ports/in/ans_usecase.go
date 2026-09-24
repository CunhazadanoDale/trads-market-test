package in

import "context"

type ANSUsecase interface {
	Import(ctx context.Context) error
}
