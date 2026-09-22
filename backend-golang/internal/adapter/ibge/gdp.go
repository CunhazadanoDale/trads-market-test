package ibge

import (
	"context"
	"fmt"
	"net/url"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
)

func (c *Client) GetGDP2023(ctx context.Context) ([]dtos.GDPSeries, error) {
	query := url.Values{}

	query.Set("localidades", "N6[all]")

	var response []dtos.GDPRecord

	err := c.Get(
		ctx,
		"/agregados/5938/periodos/2023/variaveis/37",
		query,
		&response,
	)
	if err != nil {
		return nil, fmt.Errorf("get GDP 2023: %w", err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf(
			"IBGE returned empty GDP response",
		)
	}

	if len(response[0].Resultados) == 0 {
		return nil, fmt.Errorf(
			"IBGE returned no GDP results",
		)
	}

	return response[0].Resultados[0].Series, nil
}
