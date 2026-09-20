package postgres

import (
	"context"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.StatesRepository = (*StateRepo)(nil)

type StateRepo struct {
	db *pgxpool.Pool
}

func NewStateRepo(db *pgxpool.Pool) *StateRepo {
	return &StateRepo{
		db: db,
	}
}

// Upsert implements [out.StatesRepository].
func (s *StateRepo) Upsert(ctx context.Context, state *domain.State) error {
	query := `INSERT INTO states (
			ibge_code,
			name,
			uf,
			region
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (ibge_code)
		DO UPDATE SET
			name = EXCLUDED.name,
			uf = EXCLUDED.uf,
			region = EXCLUDED.region,
			updated_at = NOW()`

	_, err := s.db.Exec(ctx, query,
		state.IBGECode,
		state.Name,
		state.UF,
		state.Region,
	)
	return err
}
