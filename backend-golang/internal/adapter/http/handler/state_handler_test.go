package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
)

var _ in.StateUseCase = (*fakeStateUseCase)(nil)

type fakeStateUseCase struct {
	states     []domain.State
	err        error
	recebeu    string
	chamouFind bool
}

func (f *fakeStateUseCase) Import(context.Context) error {
	return nil
}

func (f *fakeStateUseCase) FindAll(_ context.Context, regiao string) ([]domain.State, error) {
	f.chamouFind = true
	f.recebeu = regiao
	return f.states, f.err
}

func TestStateHandlerFindAll(t *testing.T) {
	estadosNorte := []domain.State{
		{ID: 1, IBGECode: 12, Name: "Acre", UF: "AC", Region: "Norte"},
		{ID: 2, IBGECode: 16, Name: "Amapá", UF: "AP", Region: "Norte"},
	}

	tests := []struct {
		nome         string
		query        string
		estados      []domain.State
		wantStatus   int
		wantRegiao   string
		wantChamou   bool
		wantCode     string
		wantTamanho  int
	}{
		{
			nome:        "regiao valida filtra e devolve 200",
			query:       "?regiao=Norte",
			estados:     estadosNorte,
			wantStatus:  http.StatusOK,
			wantRegiao:  "Norte",
			wantChamou:  true,
			wantTamanho: 2,
		},
		{
			nome:        "sem regiao mantem comportamento atual",
			query:       "",
			estados:     estadosNorte,
			wantStatus:  http.StatusOK,
			wantRegiao:  "",
			wantChamou:  true,
			wantTamanho: 2,
		},
		{
			nome:       "regiao invalida devolve 400 e nao chama o usecase",
			query:      "?regiao=XXX",
			wantStatus: http.StatusBadRequest,
			wantChamou: false,
			wantCode:   CodeInvalidRequest,
		},
		{
			nome:       "regiao com hifen errado devolve 400",
			query:      "?regiao=Centro%20Oeste",
			wantStatus: http.StatusBadRequest,
			wantChamou: false,
			wantCode:   CodeInvalidRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			fake := &fakeStateUseCase{states: tt.estados}
			handler := NewStateHandler(fake)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/states"+tt.query, nil)
			rec := httptest.NewRecorder()

			handler.FindAll(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, quero %d", rec.Code, tt.wantStatus)
			}

			if fake.chamouFind != tt.wantChamou {
				t.Fatalf("chamouFind = %v, quero %v", fake.chamouFind, tt.wantChamou)
			}

			if tt.wantChamou && fake.recebeu != tt.wantRegiao {
				t.Fatalf("regiao recebida = %q, quero %q", fake.recebeu, tt.wantRegiao)
			}

			if tt.wantCode != "" {
				var erro struct {
					Error struct {
						Code string `json:"code"`
					} `json:"error"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&erro); err != nil {
					t.Fatalf("erro ao decodificar resposta: %v", err)
				}
				if erro.Error.Code != tt.wantCode {
					t.Fatalf("code = %q, quero %q", erro.Error.Code, tt.wantCode)
				}
				return
			}

			var estados []dtosState
			if err := json.NewDecoder(rec.Body).Decode(&estados); err != nil {
				t.Fatalf("erro ao decodificar estados: %v", err)
			}
			if len(estados) != tt.wantTamanho {
				t.Fatalf("tamanho = %d, quero %d", len(estados), tt.wantTamanho)
			}
		})
	}
}

type dtosState struct {
	Region string `json:"region"`
}
