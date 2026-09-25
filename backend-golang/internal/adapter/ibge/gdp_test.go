package ibge

import (
	"errors"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
)

func TestGDPSeriesGDP(t *testing.T) {
	record := dtos.GDPSeries{
		Localidade: dtos.Localidade{
			ID: "1100015",
			Nivel: dtos.Nivel{
				ID:   "N6",
				Nome: "Município",
			},
			Nome: "Alta Floresta D'Oeste - RO",
		},
		Serie: map[string]string{
			"2023": "1210.60",
		},
	}

	gdp, err := record.GDP("2023")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gdp != 1210.60 {
		t.Fatalf(
			"expected GDP 1210.60, got %f",
			gdp,
		)
	}
}

func TestGDPSeriesGDPYearNotFound(t *testing.T) {
	record := dtos.GDPSeries{
		Localidade: dtos.Localidade{
			ID: "1100015",
			Nivel: dtos.Nivel{
				ID:   "N6",
				Nome: "Município",
			},
			Nome: "Alta Floresta D'Oeste - RO",
		},
		Serie: map[string]string{
			"2023": "1210.60",
		},
	}

	_, err := record.GDP("2022")
	if err == nil {
		t.Fatal("expected error when GDP year is not found")
	}
}

func TestGDPSeriesGDPInvalidValue(t *testing.T) {
	record := dtos.GDPSeries{
		Localidade: dtos.Localidade{
			ID: "1100015",
			Nivel: dtos.Nivel{
				ID:   "N6",
				Nome: "Município",
			},
			Nome: "Alta Floresta D'Oeste - RO",
		},
		Serie: map[string]string{
			"2023": "invalid",
		},
	}

	_, err := record.GDP("2023")
	if err == nil {
		t.Fatal("expected error when GDP value is invalid")
	}
}

func TestGDPSeriesGDPSuppressedValue(t *testing.T) {
	tests := []struct {
		nome  string
		serie map[string]string
		ano   string
	}{
		{
			nome:  "traco devolve valor suprimido",
			serie: map[string]string{"2023": "-"},
			ano:   "2023",
		},
		{
			nome:  "pontos devolvem valor suprimido",
			serie: map[string]string{"2023": ".."},
			ano:   "2023",
		},
		{
			nome:  "vazio devolve valor suprimido",
			serie: map[string]string{"2023": ""},
			ano:   "2023",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			record := dtos.GDPSeries{
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

			got, err := record.GDP(tt.ano)
			if !errors.Is(err, dtos.ErrSuppressedValue) {
				t.Fatalf("err = %v, quero %v", err, dtos.ErrSuppressedValue)
			}

			if got != 0 {
				t.Errorf("gdp = %f, quero 0", got)
			}
		})
	}
}
