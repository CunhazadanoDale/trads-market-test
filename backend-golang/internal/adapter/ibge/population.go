package ibge

import (
	"fmt"
	"strconv"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
)

func PopulationByYear(record dtos.PopulationRecord, year string) (int64, error) {
	value, ok := record.Serie[year]
	if !ok {
		return 0, fmt.Errorf("ano %s não encontrado na série", year)
	}

	if value == "-" || value == ".." || value == "" {
		return 0, dtos.ErrSuppressedValue
	}

	population, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("erro ao converter valor da população: %w", err)
	}

	return population, nil
}
