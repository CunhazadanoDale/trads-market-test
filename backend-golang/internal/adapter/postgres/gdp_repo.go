package postgres

import (
	"context"
	"fmt"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.GDPRepository = (*GDPRepo)(nil)

type GDPRepo struct {
	db *pgxpool.Pool
}

func NewGDPRepo(db *pgxpool.Pool) *GDPRepo {
	return &GDPRepo{
		db: db,
	}
}

func (g *GDPRepo) Upsert(ctx context.Context, ibgeCode int64, year int, gdp float64) error {
	query := `
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

	result, err := g.db.Exec(
		ctx,
		query,
		ibgeCode,
		year,
		gdp,
	)
	if err != nil {
		return fmt.Errorf(
			"upsert GDP for IBGE code %d: %w",
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
