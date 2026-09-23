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
