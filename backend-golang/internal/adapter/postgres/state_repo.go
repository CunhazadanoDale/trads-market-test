package postgres

import (
	"context"
	"fmt"

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

// FindAll implements [out.StatesRepository].
func (s *StateRepo) FindAll(ctx context.Context, regiao string) ([]domain.State, error) {
	const query = `
		SELECT
			id,
			ibge_code,
			uf,
			name,
			region
		FROM states
		WHERE ($1 = '' OR region = $1)
		ORDER BY name
	`

	rows, err := s.db.Query(ctx, query, regiao)
	if err != nil {
		return nil, fmt.Errorf("query states: %w", err)
	}
	defer rows.Close()

	states := make([]domain.State, 0)

	for rows.Next() {
		var state domain.State

		if err := rows.Scan(
			&state.ID,
			&state.IBGECode,
			&state.UF,
			&state.Name,
			&state.Region,
		); err != nil {
			return nil, fmt.Errorf("scan state: %w", err)
		}

		states = append(states, state)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate states: %w", err)
	}

	return states, nil
}
