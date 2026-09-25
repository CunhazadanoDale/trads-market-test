package ibge

import (
	"errors"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
)

func TestIncomeSeriesIncome(t *testing.T) {
	record := dtos.IncomeSeries{
		Localidade: dtos.Localidade{
			ID: "1100015",
			Nivel: dtos.Nivel{
				ID:   "N6",
				Nome: "Município",
			},
			Nome: "Alta Floresta D'Oeste - RO",
		},
		Serie: map[string]string{
			"2022": "1210.60",
		},
	}

	income, err := record.Income("2022")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if income != 1210.60 {
		t.Fatalf(
			"expected income 1210.60, got %f",
			income,
		)
	}
}

func TestIncomeSeriesIncomeYearNotFound(t *testing.T) {
	record := dtos.IncomeSeries{
		Localidade: dtos.Localidade{
			ID: "1100015",
			Nivel: dtos.Nivel{
				ID:   "N6",
				Nome: "Município",
			},
			Nome: "Alta Floresta D'Oeste - RO",
		},
		Serie: map[string]string{
			"2022": "1210.60",
		},
	}

	_, err := record.Income("2021")
	if err == nil {
		t.Fatal("expected error when income year is not found")
	}
}

func TestIncomeSeriesIncomeInvalidValue(t *testing.T) {
	record := dtos.IncomeSeries{
		Localidade: dtos.Localidade{
			ID: "1100015",
			Nivel: dtos.Nivel{
				ID:   "N6",
				Nome: "Município",
			},
			Nome: "Alta Floresta D'Oeste - RO",
		},
		Serie: map[string]string{
			"2022": "invalid",
		},
	}

	_, err := record.Income("2022")
	if err == nil {
		t.Fatal("expected error when income value is invalid")
	}
}

func TestIncomeSeriesIncomeSuppressedValue(t *testing.T) {
	tests := []struct {
		nome  string
		serie map[string]string
		ano   string
	}{
		{
			nome:  "traco devolve valor suprimido",
			serie: map[string]string{"2022": "-"},
			ano:   "2022",
		},
		{
			nome:  "pontos devolvem valor suprimido",
			serie: map[string]string{"2022": ".."},
			ano:   "2022",
		},
		{
			nome:  "vazio devolve valor suprimido",
			serie: map[string]string{"2022": ""},
			ano:   "2022",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			record := dtos.IncomeSeries{
				Localidade: dtos.Localidade{
					ID: "1100015",
					Nivel: dtos.Nivel{
						ID:   "N6",
						Nome: "Município",
					},
					Nome: "Alta Floresta D'Oeste - RO",
				},
				Serie: tt.serie,
			}

			got, err := record.Income(tt.ano)
			if !errors.Is(err, dtos.ErrSuppressedValue) {
				t.Fatalf("err = %v, quero %v", err, dtos.ErrSuppressedValue)
			}

			if got != 0 {
				t.Errorf("income = %f, quero 0", got)
			}
		})
	}
}
