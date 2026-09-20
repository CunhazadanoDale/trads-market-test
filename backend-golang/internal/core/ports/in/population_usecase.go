package in

import "context"

type PopulationUsecase interface {
	Import2022(ctx context.Context) error
}