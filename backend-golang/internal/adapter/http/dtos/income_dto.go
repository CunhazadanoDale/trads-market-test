package dtos

import (
	"fmt"
	"strconv"
)

type IncomeRecord struct {
	ID         string         `json:"id"`
	Variavel   string         `json:"variavel"`
	Unidade    string         `json:"unidade"`
	Resultados []IncomeResult `json:"resultados"`
}

type IncomeResult struct {
	Classificacoes []IncomeClassification `json:"classificacoes"`
	Series         []IncomeSeries         `json:"series"`
}

type IncomeClassification struct {
	ID        string            `json:"id"`
	Nome      string            `json:"nome"`
	Categoria map[string]string `json:"categoria"`
}

type IncomeSeries struct {
	Localidade Localidade        `json:"localidade"`
	Serie      map[string]string `json:"serie"`
}

func (r IncomeSeries) Income(year string) (float64, error) {
	value, ok := r.Serie[year]
	if !ok {
		return 0, fmt.Errorf(
			"income not found for year %s in locality %s",
			year,
			r.Localidade.Nome,
		)
	}

	income, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"parse income %q for locality %s: %w",
			value,
			r.Localidade.Nome,
			err,
		)
	}

	return income, nil
}