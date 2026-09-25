package usecases

import (
	"context"
	"fmt"
	"log"

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
		repo:       repo,
		ibgeClient: ibgeClient,
	}
}

func (c *CityUsecaseImpl) Import(ctx context.Context) error {
	states, err := c.ibgeClient.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("get states from IBGE : %w", err)
	}

	for _, state := range states {
		log.Printf("importando municipios de %s ... ", state.Sigla)

		cities, err := c.ibgeClient.GetCitiesByState(ctx, state.Sigla)
		if err != nil {
			return fmt.Errorf("capturar cidades do estado %s : %w", state.Sigla, err)
		}

		batch := make([]domain.City, 0, len(cities))

		for _, item := range cities {
			batch = append(batch, domain.City{
				IBGECode: item.ID,
				Name:     item.Nome,
			})
		}

		if err := c.repo.UpsertMany(ctx, batch, state.ID); err != nil {
			return fmt.Errorf("upsert cidades de %s : %w", state.Sigla, err)
		}

		log.Printf("%s: %d municípios importados",
			state.Sigla, len(cities))
	}

	return nil
}

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPage         = 10000
	maxPageSize     = 100
)

var ordenarPermitido = map[string]bool{
	"":          true,
	"populacao": true,
	"renda":     true,
	"pib":       true,
}

var ordemPermitida = map[string]bool{
	"":     true,
	"asc":  true,
	"desc": true,
}

func (c *CityUsecaseImpl) FindByState(
	ctx context.Context,
	stateIBGECode int64,
	filter domain.PaginacaoFilter,
) (domain.PaginacaoResponse[domain.CityWithIndicators], error) {
	filter = normalizePaginacao(filter)

	if !ordenarPermitido[filter.Ordenar] {
		return domain.PaginacaoResponse[domain.CityWithIndicators]{},
			domain.ErrInvalidOrdenar
	}

	if !ordemPermitida[filter.Ordem] {
		return domain.PaginacaoResponse[domain.CityWithIndicators]{},
			domain.ErrInvalidOrdem
	}

	cities, total, err := c.repo.FindByState(
		ctx,
		stateIBGECode,
		filter,
	)
	if err != nil {
		return domain.PaginacaoResponse[domain.CityWithIndicators]{},
			fmt.Errorf("buscar cidades do estado %d: %w", stateIBGECode, err)
	}

	if total == 0 {
		exists, err := c.repo.StateExists(ctx, stateIBGECode)
		if err != nil {
			return domain.PaginacaoResponse[domain.CityWithIndicators]{},
				fmt.Errorf("verificar estado %d: %w", stateIBGECode, err)
		}

		if !exists {
			return domain.PaginacaoResponse[domain.CityWithIndicators]{}, domain.ErrStateNotFound
		}
	}

	return domain.PaginacaoResponse[domain.CityWithIndicators]{
		Dados: cities,
		Page:  filter.Page,
		Size:  filter.Size,
		Total: total,
	}, nil
}

func (c *CityUsecaseImpl) FindByIBGECode(
	ctx context.Context,
	ibgeCode int64,
) (domain.CityDetail, error) {
	detail, err := c.repo.FindDetailByIBGECode(ctx, ibgeCode)
	if err != nil {
		return domain.CityDetail{}, fmt.Errorf("buscar cidade %d: %w", ibgeCode, err)
	}

	return *detail, nil
}

func normalizePaginacao(filter domain.PaginacaoFilter) domain.PaginacaoFilter {
	if filter.Page < defaultPage {
		filter.Page = defaultPage
	}

	if filter.Page > maxPage {
		filter.Page = maxPage
	}

	if filter.Size < 1 {
		filter.Size = defaultPageSize
	}

	if filter.Size > maxPageSize {
		filter.Size = maxPageSize
	}

	return filter
}
