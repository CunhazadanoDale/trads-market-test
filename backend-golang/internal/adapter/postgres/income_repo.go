package postgres

import (
	"context"
	"fmt"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.IncomeRepository = (*IncomeRepo)(nil)

type IncomeRepo struct {
	db *pgxpool.Pool
}

func NewIncomeRepo (db *pgxpool.Pool) *IncomeRepo {
	return &IncomeRepo{
		db: db,
	}
}

// Upsert implements [out.IncomeRepository].
func (i *IncomeRepo) Upsert(ctx context.Context, ibgeCode int64, year int, averageIncome float64) error {
	query := `INSERT INTO income_indicators (
			city_id,
			year,
			average_income
		)
		SELECT
			id,
			$2,
			$3
		FROM cities
		WHERE ibge_code = $1
		ON CONFLICT (city_id, year)
		DO UPDATE SET
			average_income = EXCLUDED.average_income,
			updated_at = NOW()
	`

	result, err := i.db.Exec(
		ctx,
		query,
		ibgeCode,
		year,
		averageIncome,
	)
	if err != nil {
		return fmt.Errorf(
			"upsert income for IBGE code %d: %w",
			ibgeCode,
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"city not found for IBGE code %d",
			ibgeCode,
		)
	}

	return nil
}
