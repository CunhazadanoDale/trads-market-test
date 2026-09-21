package ibge

import (
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
