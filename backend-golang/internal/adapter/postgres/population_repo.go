package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.PopulationRepository = (*PopulationRepo)(nil)

const populationUpsertChunk = 1000

type PopulationRepo struct {
	db *pgxpool.Pool
}

func NewPopulationRepo(db *pgxpool.Pool) *PopulationRepo {
	return &PopulationRepo{
		db: db,
	}
}

func (p *PopulationRepo) UpsertMany(ctx context.Context, rows []out.PopulationUpsert) error {
	const query = `INSERT INTO population_indicators (
			city_id,
			year,
			value,
			source
		)
		SELECT
			id,
			$2,
			$3,
			$4
		FROM cities
		WHERE ibge_code = $1
		ON CONFLICT (city_id, year)
		DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()`

	missing := 0
	for start := 0; start < len(rows); start += populationUpsertChunk {
		end := min(start+populationUpsertChunk, len(rows))

		chunkMissing, err := p.upsertChunk(ctx, query, rows[start:end])
		if err != nil {
			return err
		}

		missing += chunkMissing
	}

	if missing > 0 {
		log.Printf("population_indicators: %d indicadores sem cidade correspondente ignorados", missing)
	}

	return nil
}

func (p *PopulationRepo) upsertChunk(
	ctx context.Context,
	query string,
	chunk []out.PopulationUpsert,
) (int, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx de population_indicators: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	for _, row := range chunk {
		batch.Queue(query, row.IBGECode, row.Year, row.Value, row.Source)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	missing := 0

	for i := range chunk {
		tag, err := results.Exec()
		if err != nil {
			return 0, fmt.Errorf(
				"upsert população para o código IBGE %d: %w",
				chunk[i].IBGECode, err,
			)
		}

		if tag.RowsAffected() == 0 {
			missing++
			continue
		}
	}

	if err := results.Close(); err != nil {
		return 0, fmt.Errorf("fechar batch de population_indicators: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit de population_indicators: %w", err)
	}

	return missing, nil
}
