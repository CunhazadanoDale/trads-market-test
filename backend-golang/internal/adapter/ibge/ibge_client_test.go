package ibge

import (
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
)

func TestPopulationRecordPopulation(t *testing.T) {
	record := dtos.PopulationRecord{
		Localidade: dtos.Localidade{
			ID: "4314506",
			Nivel: dtos.Nivel{
				ID:   "N6",
				Nome: "Município",
			},
			Nome: "Pinheiro Machado - RS",
		},
		Serie: map[string]string{
			"2022": "11214",
		},
	}

	population, err := PopulationByYear(record, "2022")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if population != 11214 {
		t.Fatalf(
			"expected population 11214, got %d",
			population,
		)
	}
}



func TestPopulationRecordPopulationYearNotFound(t *testing.T) {
	record := dtos.PopulationRecord{
		Localidade: dtos.Localidade{
			ID: "4314506",
			Nivel: dtos.Nivel{
				ID:   "N6",
				Nome: "Município",
			},
			Nome: "Pinheiro Machado - RS",
		},
		Serie: map[string]string{
			"2022": "11214",
		},
	}

	_, err := PopulationByYear(record, "2021")
	if err == nil {
		t.Fatal("expected error when population year is not found")
	}
}

func TestPopulationRecordPopulationInvalidValue(t *testing.T) {
	record := dtos.PopulationRecord{
		Localidade: dtos.Localidade{
			ID: "4314506",
			Nivel: dtos.Nivel{
				ID:   "N6",
				Nome: "Município",
			},
			Nome: "Pinheiro Machado - RS",
		},
		Serie: map[string]string{
			"2022": "invalid",
		},
	}

	_, err := PopulationByYear(record, "2022")
	if err == nil {
		t.Fatal("expected error when population value is invalid")
	}
}
