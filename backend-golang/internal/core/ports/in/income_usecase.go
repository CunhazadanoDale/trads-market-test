package in

import "context"

type IncomeUsecase interface {
	Import(ctx context.Context) error
}
