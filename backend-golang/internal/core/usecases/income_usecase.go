package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
)

var _ out.IncomeRepository = (*IncomeUsecaseImpl)(nil)

const IncomeYear = 2022

type IncomeUsecaseImpl struct {
	repo       out.IncomeRepository
	ibgeClient *ibge.Client
}

func NewIncomeUsecaseImpl (repo out.IncomeRepository, ibgeClient *ibge.Client) *IncomeUsecaseImpl {
	return &IncomeUsecaseImpl{
		repo: repo,
		ibgeClient: ibgeClient,
	}
}

// Upsert implements [out.IncomeRepository].
func (i *IncomeUsecaseImpl) Upsert(ctx context.Context, ibgeCode int64, year int, averageIncome float64) error {
	records, err := i.ibgeClient.GetIncome2022(ctx)
	if err != nil {
		return fmt.Errorf(
			"get income from IBGE: %w",
			err,
		)
	}

	for _, record := range records {
		ibgeCode, err := strconv.ParseInt(
			record.Localidade.ID,
			10,
			64,
		)
		if err != nil {
			return fmt.Errorf(
				"parse IBGE code %q for locality %q: %w",
				record.Localidade.ID,
				record.Localidade.Nome,
				err,
			)
		}

		averageIncome, err := record.Income("2022")
		if err != nil {
			return fmt.Errorf(
				"get income for locality %q: %w",
				record.Localidade.Nome,
				err,
			)
		}

		if err := i.repo.Upsert(
			ctx,
			ibgeCode,
			IncomeYear,
			averageIncome,
		); err != nil {
			return fmt.Errorf(
				"persist income for locality %q: %w",
				record.Localidade.Nome,
				err,
			)
		}
	}

	return nil
}
