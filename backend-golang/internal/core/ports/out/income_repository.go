package out

import "context"

type IncomeRepository interface {
	Upsert(ctx context.Context, ibgeCode int64, year int, averageIncome float64) error
}
