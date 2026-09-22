package http

import (
	"net/http"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/handler"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(
	db *pgxpool.Pool,
	stateUseCase in.StateUseCase,
	cityUseCase in.CityUseCase,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/health/db", handler.HealthHandlerWithDBCheck(db))

	stateHandler := handler.NewStateHandler(stateUseCase)
	mux.HandleFunc("GET /api/v1/states", stateHandler.FindAll)

	cityHandler := handler.NewCityHandler(cityUseCase)
	mux.HandleFunc("GET /api/v1/states/{ibgeCode}/cities", cityHandler.FindByState)

	return mux
}
