package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
)

type fakeCityRepository struct {
	cities []domain.CityWithIndicators
	total  int

	filtroRecebido domain.PaginacaoFilter
	chamouFind     bool
}

func (f *fakeCityRepository) Upsert(context.Context, *domain.City, int64) error {
	return nil
}

func (f *fakeCityRepository) FindByState(
	_ context.Context,
	_ int64,
	filter domain.PaginacaoFilter,
) ([]domain.CityWithIndicators, int, error) {
	f.chamouFind = true
	f.filtroRecebido = filter
	return f.cities, f.total, nil
}

func (f *fakeCityRepository) StateExists(context.Context, int64) (bool, error) {
	return true, nil
}

func (f *fakeCityRepository) FindDetailByIBGECode(context.Context, int64) (*domain.CityDetail, error) {
	return nil, nil
}

func TestCityUsecaseFindByStateFiltros(t *testing.T) {
	tests := []struct {
		nome           string
		filter         domain.PaginacaoFilter
		wantErr        error
		wantChamouFind bool
		wantOrdenar    string
		wantOrdem      string
		wantNome       string
	}{
		{
			nome:           "sem ordenar e ordem usa order by de nome",
			filter:         domain.PaginacaoFilter{},
			wantChamouFind: true,
			wantOrdenar:    "",
			wantOrdem:      "",
			wantNome:       "",
		},
		{
			nome:           "ordenar e ordem validos chegam ao repo",
			filter:         domain.PaginacaoFilter{Ordenar: "renda", Ordem: "desc"},
			wantChamouFind: true,
			wantOrdenar:    "renda",
			wantOrdem:      "desc",
			wantNome:       "",
		},
		{
			nome:           "nome chega ao repo",
			filter:         domain.PaginacaoFilter{Nome: "são"},
			wantChamouFind: true,
			wantNome:       "são",
		},
		{
			nome:    "ordenar invalido retorna ErrInvalidOrdenar",
			filter:  domain.PaginacaoFilter{Ordenar: "nome"},
			wantErr: domain.ErrInvalidOrdenar,
		},
		{
			nome:    "ordem invalida retorna ErrInvalidOrdem",
			filter:  domain.PaginacaoFilter{Ordenar: "pib", Ordem: "crescente"},
			wantErr: domain.ErrInvalidOrdem,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			repo := &fakeCityRepository{}
			usecase := NewCityUsecaseImpl(repo, nil)

			_, err := usecase.FindByState(context.Background(), 12, tt.filter)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, quero %v", err, tt.wantErr)
				}
				if repo.chamouFind {
					t.Fatal("repo.FindByState não deveria ser chamado com filtro inválido")
				}
				return
			}

			if err != nil {
				t.Fatalf("err inesperado: %v", err)
			}

			if !repo.chamouFind {
				t.Fatal("repo.FindByState deveria ser chamado")
			}

			if repo.filtroRecebido.Ordenar != tt.wantOrdenar {
				t.Errorf("ordenar = %q, quero %q", repo.filtroRecebido.Ordenar, tt.wantOrdenar)
			}
			if repo.filtroRecebido.Ordem != tt.wantOrdem {
				t.Errorf("ordem = %q, quero %q", repo.filtroRecebido.Ordem, tt.wantOrdem)
			}
			if repo.filtroRecebido.Nome != tt.wantNome {
				t.Errorf("nome = %q, quero %q", repo.filtroRecebido.Nome, tt.wantNome)
			}
		})
	}
}
