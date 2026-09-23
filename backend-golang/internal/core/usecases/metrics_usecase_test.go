package usecases

import (
	"context"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
)

type fakeMetricsRepository struct {
	regiaoRecebida string
	ibgeRecebido   int64
	chamouFind     bool

	distribution domain.AgeDistribution
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
) (domain.AgeDistribution, error) {
	f.chamouFind = true
	f.regiaoRecebida = regiao
	f.ibgeRecebido = ibgeCode
	return f.distribution, nil
}

func TestMetricsUsecaseFindAgeDistributionFiltros(t *testing.T) {
	tests := []struct {
		nome       string
		regiao     string
		ibgeCode   int64
		wantRegiao string
		wantIbge   int64
	}{
		{
			nome:       "sem filtro passa zeros/vazio",
			wantRegiao: "",
			wantIbge:   0,
		},
		{
			nome:       "filtro de regiao chega ao repo",
			regiao:     "Norte",
			wantRegiao: "Norte",
			wantIbge:   0,
		},
		{
			nome:       "filtro de uf chega ao repo",
			ibgeCode:   12,
			wantRegiao: "",
			wantIbge:   12,
		},
		{
			nome:       "regiao e uf juntas chegam ao repo",
			regiao:     "Norte",
			ibgeCode:   12,
			wantRegiao: "Norte",
			wantIbge:   12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			repo := &fakeMetricsRepository{}
			usecase := NewMetricsUseCaseImpl(repo)

			_, err := usecase.FindAgeDistribution(context.Background(), tt.regiao, tt.ibgeCode)

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
		})
	}
}
