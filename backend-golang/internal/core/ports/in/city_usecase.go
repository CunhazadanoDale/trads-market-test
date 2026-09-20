package in

import "context"

type CityUseCase interface {
	Import(ctx context.Context) error
}