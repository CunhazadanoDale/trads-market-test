package usecases

import (
	"context"
	"fmt"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
)

var _ in.CityUseCase = (*CityUsecaseImpl)(nil)

type CityUsecaseImpl struct {
	repo       out.CityRepository
	ibgeClient *ibge.Client
}

func NewCityUsecaseImpl(repo out.CityRepository, ibgeClient *ibge.Client) *CityUsecaseImpl {
	return &CityUsecaseImpl{
		repo: repo,
		ibgeClient: ibgeClient,
	}
}

// Import implements [in.CityUseCase].
func (c *CityUsecaseImpl) Import(ctx context.Context) error {
	states, err := c.ibgeClient.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("get states from IBGE : %w", err)
	}

	for _, state := range states {
		fmt.Printf("importando municipios de %s ... \n", state.Sigla)

		cities, err := c.ibgeClient.GetCitiesByState(ctx, state.Sigla)
		if err != nil {
			return fmt.Errorf("capturar cidades do estado %s : %w", state.Sigla, err)
		}

		for _, item := range cities {
			city := domain.City {
				IBGECode: item.ID,
				Name: item.Nome,
			}

			if err := c.repo.Upsert(ctx, &city, state.ID); err != nil {
				return fmt.Errorf("upsert city %d (%s) : %w",
				city.IBGECode, city.Name, err)
			}
		}

		fmt.Printf("%s: %d municípios importados\n",
			state.Sigla, len(cities),)
	}

	return nil
}
