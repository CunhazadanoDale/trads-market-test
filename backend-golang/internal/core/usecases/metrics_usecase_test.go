package usecases

import (
	"context"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
)

type fakeMetricsRepository struct {
	regiaoRecebida  string
	ibgeRecebido    int64
	faixaRecebida   string
	ordenarRecebido string
	limitRecebido   int
	chamouFind      bool
	chamouFindANS   bool

	distribution domain.AgeDistribution
	ansMetrics   domain.ANSMetrics
}

func (f *fakeMetricsRepository) FindNational(context.Context) (domain.NationalMetrics, error) {
	return domain.NationalMetrics{}, nil
}

func (f *fakeMetricsRepository) FindStates(context.Context, string) ([]domain.StateMetrics, error) {
	return nil, nil
}

func (f *fakeMetricsRepository) FindTopCities(context.Context, int) (domain.TopCities, error) {
	return domain.TopCities{}, nil
}

func (f *fakeMetricsRepository) FindAgeDistribution(
	_ context.Context,
	regiao string,
	ibgeCode int64,
	faixa string,
) (domain.AgeDistribution, error) {
	f.chamouFind = true
	f.regiaoRecebida = regiao
	f.ibgeRecebido = ibgeCode
	f.faixaRecebida = faixa
	return f.distribution, nil
}

func (f *fakeMetricsRepository) FindANS(
	_ context.Context,
	regiao string,
	ibgeCode int64,
	ordenar string,
	limit int,
) (domain.ANSMetrics, error) {
	f.chamouFindANS = true
	f.regiaoRecebida = regiao
	f.ibgeRecebido = ibgeCode
	f.ordenarRecebido = ordenar
	f.limitRecebido = limit
	return f.ansMetrics, nil
}

func TestMetricsUsecaseFindAgeDistributionFiltros(t *testing.T) {
	tests := []struct {
		nome       string
		regiao     string
		ibgeCode   int64
		faixa      string
		wantRegiao string
		wantIbge   int64
		wantFaixa  string
	}{
		{
			nome:       "sem filtro passa zeros/vazio",
			wantRegiao: "",
			wantIbge:   0,
			wantFaixa:  "",
		},
		{
			nome:       "filtro de regiao chega ao repo",
			regiao:     "Norte",
			wantRegiao: "Norte",
			wantIbge:   0,
			wantFaixa:  "",
		},
		{
			nome:       "filtro de uf chega ao repo",
			ibgeCode:   12,
			wantRegiao: "",
			wantIbge:   12,
			wantFaixa:  "",
		},
		{
			nome:       "filtro de faixa chega ao repo",
			faixa:      "0 a 4 anos",
			wantRegiao: "",
			wantIbge:   0,
			wantFaixa:  "0 a 4 anos",
		},
		{
			nome:       "regiao e uf juntas chegam ao repo",
			regiao:     "Norte",
			ibgeCode:   12,
			wantRegiao: "Norte",
			wantIbge:   12,
			wantFaixa:  "",
		},
		{
			nome:       "regiao, uf e faixa juntas chegam ao repo",
			regiao:     "Norte",
			ibgeCode:   12,
			faixa:      "70 anos ou mais",
			wantRegiao: "Norte",
			wantIbge:   12,
			wantFaixa:  "70 anos ou mais",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			repo := &fakeMetricsRepository{}
			usecase := NewMetricsUseCaseImpl(repo)

			_, err := usecase.FindAgeDistribution(
				context.Background(),
				tt.regiao,
				tt.ibgeCode,
				tt.faixa,
			)

			if err != nil {
				t.Fatalf("err inesperado: %v", err)
			}

			if !repo.chamouFind {
				t.Fatal("repo.FindAgeDistribution deveria ser chamado")
			}

			if repo.regiaoRecebida != tt.wantRegiao {
				t.Errorf("regiao = %q, quero %q", repo.regiaoRecebida, tt.wantRegiao)
			}

			if repo.ibgeRecebido != tt.wantIbge {
				t.Errorf("ibge = %d, quero %d", repo.ibgeRecebido, tt.wantIbge)
			}

			if repo.faixaRecebida != tt.wantFaixa {
				t.Errorf("faixa = %q, quero %q", repo.faixaRecebida, tt.wantFaixa)
			}
		})
	}
}

func TestMetricsUsecaseFindANSFiltros(t *testing.T) {
	tests := []struct {
		nome        string
		regiao      string
		ibgeCode    int64
		ordenar     string
		limit       int
		wantRegiao  string
		wantIbge    int64
		wantOrdenar string
		wantLimit   int
	}{
		{
			nome:       "sem filtro passa zeros/vazio",
			wantRegiao: "",
			wantIbge:   0,
		},
		{
			nome:        "filtro de regiao chega ao repo",
			regiao:      "Sudeste",
			ordenar:     "beneficiarios",
			limit:       25,
			wantRegiao:  "Sudeste",
			wantIbge:    0,
			wantOrdenar: "beneficiarios",
			wantLimit:   25,
		},
		{
			nome:        "filtro de uf chega ao repo",
			ibgeCode:    35,
			ordenar:     "populacao",
			limit:       1,
			wantRegiao:  "",
			wantIbge:    35,
			wantOrdenar: "populacao",
			wantLimit:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			repo := &fakeMetricsRepository{}
			usecase := NewMetricsUseCaseImpl(repo)

			_, err := usecase.FindANS(
				context.Background(),
				tt.regiao,
				tt.ibgeCode,
				tt.ordenar,
				tt.limit,
			)

			if err != nil {
				t.Fatalf("err inesperado: %v", err)
			}

			if !repo.chamouFindANS {
				t.Fatal("repo.FindANS deveria ser chamado")
			}

			if repo.regiaoRecebida != tt.wantRegiao {
				t.Errorf("regiao = %q, quero %q", repo.regiaoRecebida, tt.wantRegiao)
			}

			if repo.ibgeRecebido != tt.wantIbge {
				t.Errorf("ibge = %d, quero %d", repo.ibgeRecebido, tt.wantIbge)
			}

			if repo.ordenarRecebido != tt.wantOrdenar {
				t.Errorf("ordenar = %q, quero %q", repo.ordenarRecebido, tt.wantOrdenar)
			}

			if repo.limitRecebido != tt.wantLimit {
				t.Errorf("limit = %d, quero %d", repo.limitRecebido, tt.wantLimit)
			}
		})
	}
}
