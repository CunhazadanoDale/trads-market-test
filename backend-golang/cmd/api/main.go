package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appRouter "github.com/CunhazadanoDale/trads-market-test/internal/adapter/http"
	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/postgres"
	"github.com/CunhazadanoDale/trads-market-test/internal/config"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/usecases"
)

func main() {
	cfg := config.LoadConfig()
	if err := cfg.ValidateAPI(); err != nil {
		log.Fatalf("configuração inválida: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := postgres.NewDBConnection(dbCtx, cfg.DatabaseUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("erro ao conectar ao database: %w", err))
	}
	defer db.Close()

	ibgeClient := ibge.NewIbgeClient(cfg.BaseUrlIBGE, cfg.BaseUrlLocalidades, http.DefaultClient)

	stateRepository := postgres.NewStateRepo(db)
	stateUseCase := usecases.NewStateUseCase(stateRepository, ibgeClient)

	cityRepository := postgres.NewCityRepository(db)
	cityUseCase := usecases.NewCityUsecaseImpl(cityRepository, ibgeClient)

	metricsRepository := postgres.NewMetricsRepository(db)
	metricsUseCase := usecases.NewMetricsUseCaseImpl(metricsRepository)

	router := appRouter.NewRouter(db, stateUseCase, cityUseCase, metricsUseCase)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	log.Printf("API rodando em %s, %s", server.Addr, cfg.AppEnv)

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	case <-ctx.Done():
		log.Println("Sinal recebido, encerrando a API")

		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("erro ao encerrar a API: %v", err)
		}

		if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}

		log.Println("API encerrada")
	}
}
