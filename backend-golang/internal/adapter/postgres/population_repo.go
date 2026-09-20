package postgres

import (
	"context"
	"fmt"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.PopulationRepository = (*PopulationRepo)(nil)

type PopulationRepo struct {
	db *pgxpool.Pool
}

func NewPopulationRepo(db *pgxpool.Pool) *PopulationRepo {
	return &PopulationRepo{
		db: db,
	}
}

// Upsert implements [out.PopulationRepository].
func (p *PopulationRepo) Upsert(ctx context.Context, ibgeCode int64, year int, value int64, source string) error {
	query := `INSERT INTO population_indicators (
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
		ON CONFLICT (city_id, year, source)
		DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()`

	result, err := p.db.Exec(ctx, query, ibgeCode, year, value, source)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cidade nao encontrada para o codigo IBGE: %d", ibgeCode)
	}

	return nil
}
