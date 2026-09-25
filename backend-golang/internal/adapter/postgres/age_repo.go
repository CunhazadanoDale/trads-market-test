package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.AgeRepository = (*AgeRepo)(nil)

const ageUpsertChunk = 1000

type AgeRepo struct {
	db *pgxpool.Pool
}

func NewAgeRepo(db *pgxpool.Pool) *AgeRepo {
	return &AgeRepo{
		db: db,
	}
}

func (a *AgeRepo) UpsertMany(
	ctx context.Context,
	rows []out.AgeUpsert,
) error {
	const query = `INSERT INTO age_indicators (
			city_id,
			year,
			age_group,
			population
		)
		SELECT
			id,
			$2,
			$3,
			$4
		FROM cities
		WHERE ibge_code = $1
		ON CONFLICT (city_id, year, age_group)
		DO UPDATE SET
			population = EXCLUDED.population,
			updated_at = NOW()
	`

	missing := 0
	for start := 0; start < len(rows); start += ageUpsertChunk {
		end := min(start+ageUpsertChunk, len(rows))

		chunkMissing, err := a.upsertChunk(ctx, query, rows[start:end])
		if err != nil {
			return err
		}

		missing += chunkMissing
	}

	if missing > 0 {
		log.Printf("age_indicators: %d indicadores sem cidade correspondente ignorados", missing)
	}

	return nil
}

func (a *AgeRepo) upsertChunk(
	ctx context.Context,
	query string,
	chunk []out.AgeUpsert,
) (int, error) {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx de age_indicators: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	for _, row := range chunk {
		batch.Queue(query, row.IBGECode, row.Year, row.AgeGroup, row.Population)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	missing := 0

	for i := range chunk {
		tag, err := results.Exec()
		if err != nil {
			return 0, fmt.Errorf(
				"upsert age %q for IBGE code %d: %w",
				chunk[i].AgeGroup, chunk[i].IBGECode, err,
			)
		}

		if tag.RowsAffected() == 0 {
			missing++
			continue
		}
	}

	if err := results.Close(); err != nil {
		return 0, fmt.Errorf("fechar batch de age_indicators: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit de age_indicators: %w", err)
	}

	return missing, nil
}
