package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
)

type fakeCityRepository struct {
	cities      []domain.CityWithIndicators
	total       int
	stateExists bool

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
	return f.stateExists, nil
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
			repo := &fakeCityRepository{stateExists: true}
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

func TestCityUsecaseFindByStateEstadoInexistente(t *testing.T) {
	repo := &fakeCityRepository{stateExists: false}
	usecase := NewCityUsecaseImpl(repo, nil)

	_, err := usecase.FindByState(context.Background(), 99, domain.PaginacaoFilter{})

	if !errors.Is(err, domain.ErrStateNotFound) {
		t.Fatalf("err = %v, quero %v", err, domain.ErrStateNotFound)
	}
}

func TestCityUsecaseFindByStateVazioComEstadoExistente(t *testing.T) {
	repo := &fakeCityRepository{stateExists: true}
	usecase := NewCityUsecaseImpl(repo, nil)

	result, err := usecase.FindByState(context.Background(), 12, domain.PaginacaoFilter{})

	if err != nil {
		t.Fatalf("err inesperada: %v", err)
	}

	if len(result.Dados) != 0 {
		t.Errorf("tamanho = %d, quero 0", len(result.Dados))
	}

	if result.Total != 0 {
		t.Errorf("total = %d, quero 0", result.Total)
	}
}

func TestNormalizePaginacao(t *testing.T) {
	tests := []struct {
		nome     string
		filter   domain.PaginacaoFilter
		wantPage int
		wantSize int
	}{
		{
			nome:     "page e size zerados viram 1 e 20",
			filter:   domain.PaginacaoFilter{},
			wantPage: 1,
			wantSize: 20,
		},
		{
			nome:     "page negativa vira 1",
			filter:   domain.PaginacaoFilter{Page: -3, Size: 10},
			wantPage: 1,
			wantSize: 10,
		},
		{
			nome:     "size negativa vira 20",
			filter:   domain.PaginacaoFilter{Page: 2, Size: -1},
			wantPage: 2,
			wantSize: 20,
		},
		{
			nome:     "size acima do maximo vira 100",
			filter:   domain.PaginacaoFilter{Page: 1, Size: 500},
			wantPage: 1,
			wantSize: 100,
		},
		{
			nome:     "valores validos passam direto",
			filter:   domain.PaginacaoFilter{Page: 3, Size: 50},
			wantPage: 3,
			wantSize: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			got := normalizePaginacao(tt.filter)

			if got.Page != tt.wantPage {
				t.Errorf("page = %d, quero %d", got.Page, tt.wantPage)
			}

			if got.Size != tt.wantSize {
				t.Errorf("size = %d, quero %d", got.Size, tt.wantSize)
			}
		})
	}
}
