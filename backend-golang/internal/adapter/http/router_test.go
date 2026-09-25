package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
)

var (
	_ in.StateUseCase   = (*fakeStateUseCase)(nil)
	_ in.CityUseCase    = (*fakeCityUseCase)(nil)
	_ in.MetricsUseCase = (*fakeMetricsUseCase)(nil)
)

type fakeStateUseCase struct {
	regiaoRecebida string
	chamouFind     bool
}

func (f *fakeStateUseCase) Import(context.Context) error {
	return nil
}

func (f *fakeStateUseCase) FindAll(_ context.Context, regiao string) ([]domain.State, error) {
	f.chamouFind = true
	f.regiaoRecebida = regiao
	return []domain.State{
		{ID: 1, IBGECode: 12, Name: "Acre", UF: "AC", Region: "Norte"},
	}, nil
}

type fakeCityUseCase struct{}

func (f *fakeCityUseCase) Import(context.Context) error {
	return nil
}

func (f *fakeCityUseCase) FindByState(
	context.Context,
	int64,
	domain.PaginacaoFilter,
) (domain.PaginacaoResponse[domain.CityWithIndicators], error) {
	return domain.PaginacaoResponse[domain.CityWithIndicators]{
		Dados: []domain.CityWithIndicators{},
	}, nil
}

func (f *fakeCityUseCase) FindByIBGECode(context.Context, int64) (domain.CityDetail, error) {
	return domain.CityDetail{}, nil
}

type fakeMetricsUseCase struct{}

func (f *fakeMetricsUseCase) FindNational(context.Context) (domain.NationalMetrics, error) {
	return domain.NationalMetrics{}, nil
}

func (f *fakeMetricsUseCase) FindStates(context.Context, string) ([]domain.StateMetrics, error) {
	return []domain.StateMetrics{}, nil
}

func (f *fakeMetricsUseCase) FindTopCities(context.Context, int) (domain.TopCities, error) {
	return domain.TopCities{}, nil
}

func (f *fakeMetricsUseCase) FindAgeDistribution(
	context.Context,
	string,
	int64,
	string,
) (domain.AgeDistribution, error) {
	return domain.AgeDistribution{}, nil
}

func (f *fakeMetricsUseCase) FindANS(
	context.Context,
	string,
	int64,
	string,
	int,
) (domain.ANSMetrics, error) {
	return domain.ANSMetrics{}, nil
}

func novoRouter(stateUseCase in.StateUseCase) http.Handler {
	return NewRouter(nil, stateUseCase, &fakeCityUseCase{}, &fakeMetricsUseCase{})
}

func TestRouterStatus(t *testing.T) {
	tests := []struct {
		nome       string
		metodo     string
		rota       string
		wantStatus int
		wantCode   string
	}{
		{
			nome:       "GET /api/v1/states devolve 200",
			metodo:     http.MethodGet,
			rota:       "/api/v1/states",
			wantStatus: http.StatusOK,
		},
		{
			nome:       "GET /api/v1/states com regiao invalida devolve 400",
			metodo:     http.MethodGet,
			rota:       "/api/v1/states?regiao=XXX",
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request",
		},
		{
			nome:       "GET em rota inexistente devolve 404 em JSON",
			metodo:     http.MethodGet,
			rota:       "/rota-inexistente",
			wantStatus: http.StatusNotFound,
			wantCode:   "not_found",
		},
		{
			nome:       "POST /api/v1/states devolve 405",
			metodo:     http.MethodPost,
			rota:       "/api/v1/states",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			nome:       "GET /health devolve 200",
			metodo:     http.MethodGet,
			rota:       "/health",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			rota := novoRouter(&fakeStateUseCase{})

			req := httptest.NewRequest(tt.metodo, tt.rota, nil)
			rec := httptest.NewRecorder()

			rota.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, quero %d", rec.Code, tt.wantStatus)
			}

			if tt.wantCode == "" {
				return
			}

			if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
				t.Fatalf("content-type = %q, quero application/json", contentType)
			}

			var payload struct {
				Error struct {
					Code string `json:"code"`
				} `json:"error"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
				t.Fatalf("erro ao decodificar resposta: %v", err)
			}
			if payload.Error.Code != tt.wantCode {
				t.Fatalf("code = %q, quero %q", payload.Error.Code, tt.wantCode)
			}
		})
	}
}

func TestRouterEstadosDevolveLista(t *testing.T) {
	rota := novoRouter(&fakeStateUseCase{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/states", nil)
	rec := httptest.NewRecorder()

	rota.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, quero %d", rec.Code, http.StatusOK)
	}

	var estados []map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&estados); err != nil {
		t.Fatalf("erro ao decodificar estados: %v", err)
	}
	if len(estados) != 1 {
		t.Fatalf("tamanho = %d, quero 1", len(estados))
	}
}

func TestRouterEstadosEncaminhaRegiao(t *testing.T) {
	tests := []struct {
		nome       string
		query      string
		wantChamou bool
		wantRegiao string
		wantStatus int
	}{
		{
			nome:       "regiao valida chega ao usecase",
			query:      "?regiao=Norte",
			wantChamou: true,
			wantRegiao: "Norte",
			wantStatus: http.StatusOK,
		},
		{
			nome:       "regiao invalida nao chama o usecase",
			query:      "?regiao=XXX",
			wantChamou: false,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			useCase := &fakeStateUseCase{}
			rota := novoRouter(useCase)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/states"+tt.query, nil)
			rec := httptest.NewRecorder()

			rota.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, quero %d", rec.Code, tt.wantStatus)
			}

			if useCase.chamouFind != tt.wantChamou {
				t.Fatalf("chamouFind = %v, quero %v", useCase.chamouFind, tt.wantChamou)
			}

			if tt.wantChamou && useCase.regiaoRecebida != tt.wantRegiao {
				t.Fatalf("regiao = %q, quero %q", useCase.regiaoRecebida, tt.wantRegiao)
			}
		})
	}
}
