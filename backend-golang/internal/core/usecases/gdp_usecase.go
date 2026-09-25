package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
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

func NewGDPUsecaseImpl(repo out.GDPRepository, ibge *ibge.Client) *GDPUsecaseImpl {
	return &GDPUsecaseImpl{
		repo: repo,
		ibge: ibge,
	}
}

func (g *GDPUsecaseImpl) Import(ctx context.Context) error {
	records, err := g.ibge.GetGDP2023(ctx)
	if err != nil {
		return fmt.Errorf(
			"get GDP from IBGE: %w",
			err,
		)
	}

	rows := make([]out.GDPUpsert, 0, len(records))
	skipped := 0

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
			if errors.Is(err, dtos.ErrSuppressedValue) {
				skipped++
				continue
			}

			return fmt.Errorf(
				"get GDP for locality %q: %w",
				record.Localidade.Nome,
				err,
			)
		}

		rows = append(rows, out.GDPUpsert{
			IBGECode: ibgeCode,
			Year:     GDPYear,
			GDP:      gdp,
		})
	}

	if err := g.repo.UpsertMany(ctx, rows); err != nil {
		return fmt.Errorf("persistir PIB: %w", err)
	}

	log.Printf("%d linhas gravadas em gdp_indicators", len(rows))
	log.Printf("%d localidades com valor suprimido ignoradas", skipped)

	return nil
}
