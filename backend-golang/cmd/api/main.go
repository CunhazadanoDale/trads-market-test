package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	appRouter "github.com/CunhazadanoDale/trads-market-test/internal/adapter/http"
	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/postgres"
	"github.com/CunhazadanoDale/trads-market-test/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := postgres.NewDBConnection(context, cfg.DatabaseUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("erro ao conectar ao database: %w", err))
	}
	defer db.Close()

	router := appRouter.NewRouter(db)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	log.Printf("API rodando em %s, %s", server.Addr, cfg.AppEnv)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
