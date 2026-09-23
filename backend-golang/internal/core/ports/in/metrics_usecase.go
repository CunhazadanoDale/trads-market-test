package in

import (
	"context"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
)

type MetricsUseCase interface {
	FindNational(ctx context.Context) (domain.NationalMetrics, error)
	FindStates(ctx context.Context, regiao string) ([]domain.StateMetrics, error)
	FindTopCities(ctx context.Context, limit int) (domain.TopCities, error)
	FindAgeDistribution(ctx context.Context, regiao string, ibgeCode int64) (domain.AgeDistribution, error)
}
