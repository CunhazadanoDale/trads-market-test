package main

import (
	"log"
	"net/http"

	appRouter "github.com/CunhazadanoDale/trads-market-test/internal/adapter/http"
	"github.com/CunhazadanoDale/trads-market-test/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	router := appRouter.NewRouter()

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	log.Printf("API rodando em %s, %s", server.Addr, cfg.AppEnv)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
