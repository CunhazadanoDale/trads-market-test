package postgres

import (
	"context"
	"fmt"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.ANSRepository = (*ANSRepo)(nil)

const ansUpsertChunk = 1000

const ansUpsertQuery = `
	WITH target AS (
		SELECT ibge_code
		FROM cities
		WHERE ibge_code / 10 = $1
		LIMIT 1
	)
	INSERT INTO ans_beneficiarios (ibge_code, ano, beneficiarios_ativos, fonte)
	SELECT ibge_code, $2, $3, $4
	FROM target
	ON CONFLICT (ibge_code, ano)
	DO UPDATE SET
		beneficiarios_ativos = EXCLUDED.beneficiarios_ativos,
		fonte = EXCLUDED.fonte,
		updated_at = NOW()
`

type ANSRepo struct {
	db *pgxpool.Pool
}

func NewANSRepo(db *pgxpool.Pool) *ANSRepo {
	return &ANSRepo{
		db: db,
	}
}

func (a *ANSRepo) UpsertMany(
	ctx context.Context,
	rows []out.ANSUpsert,
) ([]int64, error) {
	missing := make([]int64, 0)

	for start := 0; start < len(rows); start += ansUpsertChunk {
		end := min(start+ansUpsertChunk, len(rows))

		codes, err := a.upsertChunk(ctx, rows[start:end])
		if err != nil {
			return nil, err
		}

		missing = append(missing, codes...)
	}

	return missing, nil
}

func (a *ANSRepo) upsertChunk(
	ctx context.Context,
	chunk []out.ANSUpsert,
) ([]int64, error) {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx de ans_beneficiarios: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	for _, row := range chunk {
		batch.Queue(ansUpsertQuery, row.IBGECode, row.Year, row.Beneficiaries, row.Source)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	missing := make([]int64, 0)
	for i := range chunk {
		tag, err := results.Exec()
		if err != nil {
			return nil, fmt.Errorf("upsert ans para o código %d: %w", chunk[i].IBGECode, err)
		}

		if tag.RowsAffected() == 0 {
			missing = append(missing, chunk[i].IBGECode)
		}
	}

	if err := results.Close(); err != nil {
		return nil, fmt.Errorf("fechar batch de ans_beneficiarios: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit de ans_beneficiarios: %w", err)
	}

	return missing, nil
}
