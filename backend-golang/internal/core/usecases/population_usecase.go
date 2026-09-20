package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
)

var _ in.PopulationUsecase = (*PopulationUsecaseImpl)(nil)

const (
	PopulationYear         = 2022
	PopulationYearToString = "2022"
	PopulationSource       = "CENSO_2022"
)

type PopulationUsecaseImpl struct {
	repo       out.PopulationRepository
	ibgeClient *ibge.Client
}

func NewPopulationUsecaseImpl(repo out.PopulationRepository, ibgeClient *ibge.Client) *PopulationUsecaseImpl {
	return &PopulationUsecaseImpl{
		repo:       repo,
		ibgeClient: ibgeClient,
	}
}

// Import2022 implements [in.PopulationUsecase].
func (p *PopulationUsecaseImpl) Import2022(ctx context.Context) error {
	records, err := p.ibgeClient.GetPopulation2022(ctx)
	if err != nil {
		return fmt.Errorf("get population from IBGE: %w", err)
	}

	for _, record := range records {
		ibgeCode, err := strconv.ParseInt(record.Localidade.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("parse IBGE code %q for locality %q: %w", record.Localidade.ID,
				record.Localidade.Nome, err)
		}

		population, err := ibge.PopulationByYear(record, PopulationYearToString)
		if err != nil {
			return fmt.Errorf("get population for locality %q: %w",
				record.Localidade.Nome, err)
		}

		if err := p.repo.Upsert(ctx, ibgeCode, PopulationYear, population, PopulationSource); err != nil {
			return fmt.Errorf("persist population for locality %q: %w",
				record.Localidade.Nome, err)
		}
	}

	return nil
}
