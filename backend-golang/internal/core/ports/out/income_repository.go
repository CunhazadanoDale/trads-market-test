package out

import "context"

type IncomeUpsert struct {
	IBGECode      int64
	Year          int
	AverageIncome float64
}

type IncomeRepository interface {
	UpsertMany(ctx context.Context, rows []IncomeUpsert) error
}
