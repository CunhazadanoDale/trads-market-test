package postgres

import (
	"context"
	"fmt"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.GDPRepository = (*GDPRepo)(nil)

const gdpUpsertChunk = 1000

type GDPRepo struct {
	db *pgxpool.Pool
}

func NewGDPRepo(db *pgxpool.Pool) *GDPRepo {
	return &GDPRepo{
		db: db,
	}
}

func (g *GDPRepo) UpsertMany(ctx context.Context, rows []out.GDPUpsert) error {
	const query = `
		INSERT INTO gdp_indicators (
			city_id,
			year,
			gdp
		)
		SELECT
			id,
			$2,
			$3
		FROM cities
		WHERE ibge_code = $1
		ON CONFLICT (city_id, year)
		DO UPDATE SET
			gdp = EXCLUDED.gdp,
			updated_at = NOW()
	`

	for start := 0; start < len(rows); start += gdpUpsertChunk {
		end := min(start+gdpUpsertChunk, len(rows))

		if err := g.upsertChunk(ctx, query, rows[start:end]); err != nil {
			return err
		}
	}

	return nil
}

func (g *GDPRepo) upsertChunk(
	ctx context.Context,
	query string,
	chunk []out.GDPUpsert,
) error {
	tx, err := g.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx de gdp_indicators: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	for _, row := range chunk {
		batch.Queue(query, row.IBGECode, row.Year, row.GDP)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for i := range chunk {
		tag, err := results.Exec()
		if err != nil {
			return fmt.Errorf(
				"upsert PIB para o código IBGE %d: %w",
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
		return fmt.Errorf("fechar batch de gdp_indicators: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit de gdp_indicators: %w", err)
	}

	return nil
}
