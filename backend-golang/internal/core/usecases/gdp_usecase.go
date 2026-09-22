package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
)

var _ in.GDPUsecase = (*GDPUsecaseImpl)(nil)

const GDPYear = 2023

type GDPUsecaseImpl struct {
	repo out.GDPRepository
	ibge *ibge.Client
}

func NewGDPUsecaseImpl (repo out.GDPRepository, ibge *ibge.Client) *GDPUsecaseImpl {
	return &GDPUsecaseImpl{
		repo: repo,
		ibge: ibge,
	}
}

// Import implements [in.GDPUsecase].
func (g *GDPUsecaseImpl) Import(ctx context.Context) error {
	records, err := g.ibge.GetGDP2023(ctx)
	if err != nil {
		return fmt.Errorf(
			"get GDP from IBGE: %w",
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

		gdp, err := record.GDP("2023")
		if err != nil {
			return fmt.Errorf(
				"get GDP for locality %q: %w",
				record.Localidade.Nome,
				err,
			)
		}

		if err := g.repo.Upsert(
			ctx,
			ibgeCode,
			GDPYear,
			gdp,
		); err != nil {
			return fmt.Errorf(
				"persist GDP for locality %q: %w",
				record.Localidade.Nome,
				err,
			)
		}
	}

	return nil
}
