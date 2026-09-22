package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
)

type CityHandler struct {
	useCase in.CityUseCase
}

func NewCityHandler(
	useCase in.CityUseCase,
) *CityHandler {
	return &CityHandler{
		useCase: useCase,
	}
}

func (h *CityHandler) FindByState(
	w http.ResponseWriter,
	r *http.Request,
) {
	stateIBGECode, err := strconv.ParseInt(
		r.PathValue("ibgeCode"),
		10,
		64,
	)
	if err != nil {
		http.Error(
			w,
			"invalid state IBGE code",
			http.StatusBadRequest,
		)
		return
	}

	filter := domain.PaginacaoFilter{
		Page: parseQueryInt(r, "page", 1),
		Size: parseQueryInt(r, "pageSize", 20),
	}

	result, err := h.useCase.FindByState(
		r.Context(),
		stateIBGECode,
		filter,
	)
	if err != nil {
		http.Error(
			w,
			"failed to find cities",
			http.StatusInternalServerError,
		)
		return
	}

	response := domain.PaginacaoResponse[dtos.CityResponse]{
		Dados: make([]dtos.CityResponse, 0, len(result.Dados)),
		Page:  result.Page,
		Size:  result.Size,
		Total: result.Total,
	}

	for _, city := range result.Dados {
		response.Dados = append(response.Dados, dtos.NewCityResponse(city))
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func parseQueryInt(
	r *http.Request,
	name string,
	defaultValue int,
) int {
	value := r.URL.Query().Get(name)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return result
}
