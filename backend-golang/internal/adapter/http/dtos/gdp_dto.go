package dtos

import (
	"fmt"
	"strconv"
)

type GDPRecord struct {
	ID         string      `json:"id"`
	Variavel   string      `json:"variavel"`
	Unidade    string      `json:"unidade"`
	Resultados []GDPResult `json:"resultados"`
}

type GDPResult struct {
	Classificacoes []GDPClassification `json:"classificacoes"`
	Series         []GDPSeries         `json:"series"`
}

type GDPClassification struct {
	ID        string            `json:"id"`
	Nome      string            `json:"nome"`
	Categoria map[string]string `json:"categoria"`
}

type GDPSeries struct {
	Localidade Localidade        `json:"localidade"`
	Serie      map[string]string `json:"serie"`
}

func (r GDPSeries) GDP(year string) (float64, error) {
	value, ok := r.Serie[year]
	if !ok {
		return 0, fmt.Errorf(
			"GDP not found for year %s in locality %s",
			year,
			r.Localidade.Nome,
		)
	}

	if value == "-" || value == ".." || value == "" {
		return 0, ErrSuppressedValue
	}

	gdp, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"parse GDP %q for locality %s: %w",
			value,
			r.Localidade.Nome,
			err,
		)
	}

	return gdp, nil
}
