package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
)

var _ in.MetricsUseCase = (*fakeMetricsUseCase)(nil)

type fakeMetricsUseCase struct {
	topErr    error
	topLimit  int
	chamouTop bool

	ansChamou bool
	ansRegiao string
	ansIbge   int64
	ansErr    error
}

func (f *fakeMetricsUseCase) FindNational(context.Context) (domain.NationalMetrics, error) {
	return domain.NationalMetrics{}, nil
}

func (f *fakeMetricsUseCase) FindStates(context.Context, string) ([]domain.StateMetrics, error) {
	return nil, nil
}

func (f *fakeMetricsUseCase) FindAgeDistribution(
	context.Context,
	string,
	int64,
) (domain.AgeDistribution, error) {
	return domain.AgeDistribution{}, nil
}

func (f *fakeMetricsUseCase) FindTopCities(_ context.Context, limit int) (domain.TopCities, error) {
	f.chamouTop = true
	f.topLimit = limit
	return domain.TopCities{}, f.topErr
}

func (f *fakeMetricsUseCase) FindANS(
	_ context.Context,
	regiao string,
	ibgeCode int64,
) (domain.ANSMetrics, error) {
	f.ansChamou = true
	f.ansRegiao = regiao
	f.ansIbge = ibgeCode
	return domain.ANSMetrics{}, f.ansErr
}

func TestMetricsHandlerFindTopCities(t *testing.T) {
	tests := []struct {
		nome       string
		query      string
		topErr     error
		wantStatus int
		wantCode   string
		wantChamou bool
		wantLimit  int
	}{
		{
			nome:       "limit ausente usa 10",
			query:      "",
			wantStatus: http.StatusOK,
			wantChamou: true,
			wantLimit:  10,
		},
		{
			nome:       "limit valido chega ao usecase",
			query:      "?limit=50",
			wantStatus: http.StatusOK,
			wantChamou: true,
			wantLimit:  50,
		},
		{
			nome:       "limit no minimo",
			query:      "?limit=1",
			wantStatus: http.StatusOK,
			wantChamou: true,
			wantLimit:  1,
		},
		{
			nome:       "limit no maximo",
			query:      "?limit=100",
			wantStatus: http.StatusOK,
			wantChamou: true,
			wantLimit:  100,
		},
		{
			nome:       "limit nao numerico devolve 400",
			query:      "?limit=abc",
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidRequest,
			wantChamou: false,
		},
		{
			nome:       "limit zero devolve 400",
			query:      "?limit=0",
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidRequest,
			wantChamou: false,
		},
		{
			nome:       "limit acima de 100 devolve 400",
			query:      "?limit=101",
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidRequest,
			wantChamou: false,
		},
		{
			nome:       "erro do usecase devolve 500 em JSON",
			query:      "?limit=10",
			topErr:     errors.New("banco indisponivel"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   CodeInternalError,
			wantChamou: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			fake := &fakeMetricsUseCase{topErr: tt.topErr}
			handler := NewMetricsHandler(fake)

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/dashboard/top-cities"+tt.query,
				nil,
			)
			rec := httptest.NewRecorder()

			handler.FindTopCities(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, quero %d", rec.Code, tt.wantStatus)
			}

			if fake.chamouTop != tt.wantChamou {
				t.Fatalf("chamouTop = %v, quero %v", fake.chamouTop, tt.wantChamou)
			}

			if tt.wantStatus == http.StatusOK && fake.topLimit != tt.wantLimit {
				t.Fatalf("limit = %d, quero %d", fake.topLimit, tt.wantLimit)
			}

			if tt.wantCode != "" {
				var body struct {
					Error struct {
						Code string `json:"code"`
					} `json:"error"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
					t.Fatalf("erro ao decodificar resposta: %v", err)
				}
				if body.Error.Code != tt.wantCode {
					t.Fatalf("code = %q, quero %q", body.Error.Code, tt.wantCode)
				}
			}
		})
	}
}

func TestMetricsHandlerFindANS(t *testing.T) {
	tests := []struct {
		nome       string
		query      string
		ansErr     error
		wantStatus int
		wantCode   string
		wantChamou bool
		wantRegiao string
		wantIbge   int64
	}{
		{
			nome:       "sem filtro devolve 200",
			query:      "",
			wantStatus: http.StatusOK,
			wantChamou: true,
		},
		{
			nome:       "regiao valida chega ao usecase",
			query:      "?regiao=Sul",
			wantStatus: http.StatusOK,
			wantChamou: true,
			wantRegiao: "Sul",
		},
		{
			nome:       "regiao invalida devolve 400",
			query:      "?regiao=Leste",
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidRequest,
			wantChamou: false,
		},
		{
			nome:       "ibge valido chega ao usecase",
			query:      "?ibge=35",
			wantStatus: http.StatusOK,
			wantChamou: true,
			wantIbge:   35,
		},
		{
			nome:       "ibge nao numerico devolve 400",
			query:      "?ibge=abc",
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidRequest,
			wantChamou: false,
		},
		{
			nome:       "ibge zero devolve 400",
			query:      "?ibge=0",
			wantStatus: http.StatusBadRequest,
			wantCode:   CodeInvalidRequest,
			wantChamou: false,
		},
		{
			nome:       "estado inexistente devolve 404",
			query:      "?ibge=99",
			ansErr:     domain.ErrStateNotFound,
			wantStatus: http.StatusNotFound,
			wantCode:   CodeStateNotFound,
			wantChamou: true,
			wantIbge:   99,
		},
		{
			nome:       "erro do usecase devolve 500",
			query:      "",
			ansErr:     errors.New("banco indisponivel"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   CodeInternalError,
			wantChamou: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			fake := &fakeMetricsUseCase{ansErr: tt.ansErr}
			handler := NewMetricsHandler(fake)

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/dashboard/ans"+tt.query,
				nil,
			)
			rec := httptest.NewRecorder()

			handler.FindANS(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, quero %d", rec.Code, tt.wantStatus)
			}

			if fake.ansChamou != tt.wantChamou {
				t.Fatalf("ansChamou = %v, quero %v", fake.ansChamou, tt.wantChamou)
			}

			if tt.wantChamou && fake.ansRegiao != tt.wantRegiao {
				t.Fatalf("regiao = %q, quero %q", fake.ansRegiao, tt.wantRegiao)
			}

			if tt.wantChamou && fake.ansIbge != tt.wantIbge {
				t.Fatalf("ibge = %d, quero %d", fake.ansIbge, tt.wantIbge)
			}

			if tt.wantCode != "" {
				var body struct {
					Error struct {
						Code string `json:"code"`
					} `json:"error"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
					t.Fatalf("erro ao decodificar resposta: %v", err)
				}
				if body.Error.Code != tt.wantCode {
					t.Fatalf("code = %q, quero %q", body.Error.Code, tt.wantCode)
				}
			}
		})
	}
}
