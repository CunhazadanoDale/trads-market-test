package dtos

import "github.com/CunhazadanoDale/trads-market-test/internal/core/domain"

type CityResponse struct {
	ID       int64  `json:"id"`
	IBGECode int64  `json:"ibge_code"`
	Name     string `json:"name"`
}

func NewCityResponse(city domain.City) CityResponse {
	return CityResponse{
		ID:       city.ID,
		IBGECode: city.IBGECode,
		Name:     city.Name,
	}
}
