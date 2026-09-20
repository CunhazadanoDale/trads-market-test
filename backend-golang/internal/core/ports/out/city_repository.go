package out

import (
	"context"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
)

type CityRepository interface {
	Upsert(ctx context.Context, city *domain.City, stateIBGECode int64) error
}