package http

import (
	"net/http"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/handler"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, stateQueryUseCase in.StateUseCase) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/health/db", handler.HealthHandlerWithDBCheck(db))

	stateHandler := handler.NewStateHandler(stateQueryUseCase)
	mux.HandleFunc("GET /api/v1/states", stateHandler.FindAll)

	return mux
}
