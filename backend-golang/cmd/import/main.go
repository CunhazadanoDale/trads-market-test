package main

import (
	"context"
	"log"
	"time"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/postgres"
	"github.com/CunhazadanoDale/trads-market-test/internal/config"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/usecases"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		12*time.Minute,
	)
	defer cancel()

	db, err := postgres.NewDBConnection(ctx, cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("falha ao conectar com o database: %v", err)
	}
	defer db.Close()

	ibgeClient := ibge.NewIbgeClient(cfg.BaseUrlIBGE, cfg.BaseUrlLocalidades, nil)

	importStart := time.Now()

	stateRepository := postgres.NewStateRepo(db)
	stateService := usecases.NewStateUseCase(stateRepository, ibgeClient)

	stageStart := time.Now()
	if err := stateService.Import(ctx); err != nil {
		log.Fatalf("falha ao importar estados: %v", err)
	}

	log.Printf("etapa estados: %s", time.Since(stageStart))

	cityRepository := postgres.NewCityRepository(db)
	cityService := usecases.NewCityUsecaseImpl(cityRepository, ibgeClient)

	stageStart = time.Now()
	if err := cityService.Import(ctx); err != nil {
		log.Fatal("falha ao importar cidades: %w", err)
	}

	log.Printf("etapa cidades: %s", time.Since(stageStart))

	populationRepository := postgres.NewPopulationRepo(db)
	populationService := usecases.NewPopulationUsecaseImpl(populationRepository, ibgeClient)

	stageStart = time.Now()
	if err := populationService.Import2022(ctx); err != nil {
		log.Fatalf("falha ao importar população: %v", err)
	}

	log.Printf("etapa população: %s", time.Since(stageStart))

	incomeRepository := postgres.NewIncomeRepo(db)
	incomeService := usecases.NewIncomeUsecaseImpl(incomeRepository, ibgeClient)

	stageStart = time.Now()
	if err := incomeService.Import(ctx); err != nil {
		log.Fatalf("falha ao importar renda: %v", err)
	}

	log.Printf("etapa renda: %s", time.Since(stageStart))

	gdpRepository := postgres.NewGDPRepo(db)
	gdpUseCase := usecases.NewGDPUsecaseImpl(
		gdpRepository,
		ibgeClient,
	)

	stageStart = time.Now()
	if err := gdpUseCase.Import(ctx); err != nil {
		log.Fatalf("falha ao importar PIB: %v", err)
	}

	log.Printf("etapa PIB: %s", time.Since(stageStart))

	ageRepository := postgres.NewAgeRepo(db)
	ageUsecase := usecases.NewAgeUsecaseImpl(ageRepository, ibgeClient)

	stageStart = time.Now()
	if err := ageUsecase.Import(ctx); err != nil {
		log.Fatalf("falha ao importar faixa etária: %v", err)
	}

	log.Printf("etapa faixa etária: %s", time.Since(stageStart))

	log.Printf("IBGE import successfully em %s", time.Since(importStart))
}