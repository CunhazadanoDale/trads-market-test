package handler

import (
	"encoding/json"
	"net/http"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
)

type StateHandler struct {
	useCase in.StateUseCase
}

func NewStateHandler(
	useCase in.StateUseCase,
) *StateHandler {
	return &StateHandler{
		useCase: useCase,
	}
}

func (h *StateHandler) FindAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	states, err := h.useCase.FindAll(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to find states",
			http.StatusInternalServerError,
		)
		return
	}

	response := make([]dtos.StateResponse, 0, len(states))
	for _, state := range states {
		response = append(response, dtos.NewStateResponse(state))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
