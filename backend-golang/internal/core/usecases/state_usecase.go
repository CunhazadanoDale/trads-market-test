package usecases

import (
	"context"
	"fmt"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
)

var _ in.StateUseCase = (*StateUseCaseImpl)(nil)

type StateUseCaseImpl struct {
	repo       out.StatesRepository
	ibgeClient *ibge.Client
}

func NewStateUseCase(repo out.StatesRepository, ibgeClient *ibge.Client) *StateUseCaseImpl {
	return &StateUseCaseImpl{
		repo:       repo,
		ibgeClient: ibgeClient,
	}
}

// Import implements [in.StateUseCase].
func (s *StateUseCaseImpl) Import(ctx context.Context) error {
	states, err := s.ibgeClient.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("sem captura de estados pelo IBGE: %w", err)
	}

	for _, item := range states {
		state := domain.State{
			IBGECode: item.ID,
			Name: item.Nome,
			UF: item.Sigla,
			Region: item.Regiao.Nome,
		}

		if err := s.repo.Upsert(ctx, &state); err != nil {
			return fmt.Errorf("falha ao inserir/atualizar estado %s: %w", state.Name, err)
		}
	}

	return nil
}

// FindAll implements [in.StateQueryUseCase].
func (s *StateUseCaseImpl) FindAll(ctx context.Context) ([]domain.State, error) {
	states, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("buscar estados: %w", err)
	}

	return states, nil
}
