package usecases

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
)

var _ in.IncomeUsecase = (*IncomeUsecaseImpl)(nil)

const IncomeYear = 2022

type IncomeUsecaseImpl struct {
	repo       out.IncomeRepository
	ibgeClient *ibge.Client
}

func NewIncomeUsecaseImpl(repo out.IncomeRepository, ibgeClient *ibge.Client) *IncomeUsecaseImpl {
	return &IncomeUsecaseImpl{
		repo:       repo,
		ibgeClient: ibgeClient,
	}
}

func (i *IncomeUsecaseImpl) Import(ctx context.Context) error {
	records, err := i.ibgeClient.GetIncome2022(ctx)
	if err != nil {
		return fmt.Errorf(
			"get income from IBGE: %w",
			err,
		)
	}

	rows := make([]out.IncomeUpsert, 0, len(records))

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

		rows = append(rows, out.IncomeUpsert{
			IBGECode:      ibgeCode,
			Year:          IncomeYear,
			AverageIncome: averageIncome,
		})
	}

	if err := i.repo.UpsertMany(ctx, rows); err != nil {
		return fmt.Errorf("persistir renda: %w", err)
	}

	log.Printf("%d linhas gravadas em income_indicators", len(rows))

	return nil
}
