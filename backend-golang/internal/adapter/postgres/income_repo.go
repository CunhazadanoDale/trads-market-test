package postgres

import (
	"context"
	"fmt"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.IncomeRepository = (*IncomeRepo)(nil)

const incomeUpsertChunk = 1000

type IncomeRepo struct {
	db *pgxpool.Pool
}

func NewIncomeRepo(db *pgxpool.Pool) *IncomeRepo {
	return &IncomeRepo{
		db: db,
	}
}

func (i *IncomeRepo) UpsertMany(ctx context.Context, rows []out.IncomeUpsert) error {
	const query = `INSERT INTO income_indicators (
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

	for start := 0; start < len(rows); start += incomeUpsertChunk {
		end := min(start+incomeUpsertChunk, len(rows))

		if err := i.upsertChunk(ctx, query, rows[start:end]); err != nil {
			return err
		}
	}

	return nil
}

func (i *IncomeRepo) upsertChunk(
	ctx context.Context,
	query string,
	chunk []out.IncomeUpsert,
) error {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx de income_indicators: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	for _, row := range chunk {
		batch.Queue(query, row.IBGECode, row.Year, row.AverageIncome)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for i := range chunk {
		tag, err := results.Exec()
		if err != nil {
			return fmt.Errorf(
				"upsert renda para o código IBGE %d: %w",
				chunk[i].IBGECode, err,
			)
		}

		if tag.RowsAffected() == 0 {
			return fmt.Errorf(
				"city not found for IBGE code %d",
				chunk[i].IBGECode,
			)
		}
	}

	if err := results.Close(); err != nil {
		return fmt.Errorf("fechar batch de income_indicators: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit de income_indicators: %w", err)
	}

	return nil
}
