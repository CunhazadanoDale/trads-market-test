package out

import (
	"context"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
)

type StatesRepository interface {
	Upsert(ctx context.Context, state *domain.State) error
}