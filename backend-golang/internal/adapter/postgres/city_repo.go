package postgres

import (
	"context"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.CityRepository = (*CityRepo)(nil)

type CityRepo struct {
	db *pgxpool.Pool
}

func NewCityRepository(db *pgxpool.Pool) *CityRepo {
	return &CityRepo{
		db: db,
	}
}

// Upsert implements [out.CityRepository].
func (c *CityRepo) Upsert(ctx context.Context, city *domain.City, stateIBGECode int64) error {
	query := `INSERT INTO cities (
			ibge_code,
			state_id,
			name
		)
		SELECT
			$1,
			id,
			$2
		FROM states
		WHERE ibge_code = $3
		ON CONFLICT (ibge_code)
		DO UPDATE SET
			state_id = EXCLUDED.state_id,
			name = EXCLUDED.name,
			updated_at = NOW()
	`

	_, err := c.db.Exec(
		ctx, query, city.IBGECode, city.Name, stateIBGECode,
	)

	return err
}
