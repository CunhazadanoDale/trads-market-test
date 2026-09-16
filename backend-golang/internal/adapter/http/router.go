package http

import (
	"net/http"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/handler"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)

	return mux
}