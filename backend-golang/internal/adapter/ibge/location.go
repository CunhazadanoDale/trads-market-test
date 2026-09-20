package ibge

import (
	"context"
	"net/http"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
)

func (c *Client) GetStates(ctx context.Context) ([]dtos.StateRecord, error) {
	var response []dtos.StateRecord
	err := c.GetFromBase(ctx, "/localidades/estados", nil, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) GetCitiesByState(ctx context.Context, stateUF string) ([]dtos.CityRecord, error) {
	var response []dtos.CityRecord
	err := c.GetFromBase(ctx, "/localidades/estados/"+stateUF+"/municipios", nil, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var _ = http.MethodGet