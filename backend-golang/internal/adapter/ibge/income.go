package ibge

import (
	"context"
	"fmt"
	"net/url"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
)

func (c *Client) GetIncome2022(ctx context.Context) ([]dtos.IncomeSeries, error) {
	query := url.Values{}

	query.Set("localidades", "N6[all]")
	query.Add("classificacao", "2[6794]")
	query.Add("classificacao", "86[95251]")
	query.Add("classificacao", "58[95253]")

	var response []dtos.IncomeRecord

	err := c.Get(ctx,
		"/agregados/10295/periodos/2022/variaveis/13431",
		query,
		&response)
	if err != nil {
		return nil, fmt.Errorf("get income 2022: %w", err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf(
			"IBGE returned empty income response",
		)
	}

	if len(response[0].Resultados) == 0 {
		return nil, fmt.Errorf(
			"IBGE returned no income results",
		)
	}

	return response[0].Resultados[0].Series, nil
}
