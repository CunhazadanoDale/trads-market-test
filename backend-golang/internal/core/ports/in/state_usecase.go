package in

import (
	"context"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
)

type StateUseCase interface {
	Import(ctx context.Context) error
	FindAll(ctx context.Context, regiao string) ([]domain.State, error)
}
