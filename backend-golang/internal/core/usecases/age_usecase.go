package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
)

var _ in.AgeUsecase = (*AgeUsecaseImpl)(nil)

const AgeYear = 2022

type AgeUsecaseImpl struct {
	repo       out.AgeRepository
	ibgeClient *ibge.Client
}

func NewAgeUsecaseImpl(
	repo out.AgeRepository,
	ibgeClient *ibge.Client,
) *AgeUsecaseImpl {
	return &AgeUsecaseImpl{
		repo:       repo,
		ibgeClient: ibgeClient,
	}
}

func (a *AgeUsecaseImpl) Import(ctx context.Context) error {
	entries, err := a.ibgeClient.GetAgeGroups2022(ctx)
	if err != nil {
		return fmt.Errorf(
			"get age groups from IBGE: %w",
			err,
		)
	}

	rows := make([]out.AgeUpsert, 0, len(entries))

	for _, entry := range entries {
		ibgeCode, err := strconv.ParseInt(
			entry.Localidade.ID,
			10,
			64,
		)
		if err != nil {
			return fmt.Errorf(
				"parse IBGE code %q for locality %q: %w",
				entry.Localidade.ID,
				entry.Localidade.Nome,
				err,
			)
		}

		population, err := entry.Population("2022")
		if err != nil {
			return fmt.Errorf(
				"get population for locality %q: %w",
				entry.Localidade.Nome,
				err,
			)
		}

		rows = append(rows, out.AgeUpsert{
			IBGECode:   ibgeCode,
			Year:       AgeYear,
			AgeGroup:   entry.AgeGroup,
			Population: population,
		})
	}

	if err := a.repo.UpsertMany(ctx, rows); err != nil {
		return fmt.Errorf("persist age rows: %w", err)
	}

	return nil
}
