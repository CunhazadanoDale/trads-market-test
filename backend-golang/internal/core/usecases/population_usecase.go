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

func (p *PopulationUsecaseImpl) Import2022(ctx context.Context) error {
	records, err := p.ibgeClient.GetPopulation2022(ctx)
	if err != nil {
		return fmt.Errorf("get population from IBGE: %w", err)
	}

	rows := make([]out.PopulationUpsert, 0, len(records))

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

		rows = append(rows, out.PopulationUpsert{
			IBGECode: ibgeCode,
			Year:     PopulationYear,
			Value:    population,
			Source:   PopulationSource,
		})
	}

	if err := p.repo.UpsertMany(ctx, rows); err != nil {
		return fmt.Errorf("persistir população: %w", err)
	}

	log.Printf("%d linhas gravadas em population_indicators", len(rows))

	return nil
}
