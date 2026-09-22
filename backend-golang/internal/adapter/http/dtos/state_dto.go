package dtos

import "github.com/CunhazadanoDale/trads-market-test/internal/core/domain"

type StateResponse struct {
	ID       int64  `json:"id"`
	IBGECode int64  `json:"ibge_code"`
	Name     string `json:"name"`
	UF       string `json:"uf"`
	Region   string `json:"region"`
}

func NewStateResponse(state domain.State) StateResponse {
	return StateResponse{
		ID:       state.ID,
		IBGECode: state.IBGECode,
		Name:     state.Name,
		UF:       state.UF,
		Region:   state.Region,
	}
}
