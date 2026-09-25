package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ans"
	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/postgres"
	"github.com/CunhazadanoDale/trads-market-test/internal/config"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	target := "ibge"
	if len(os.Args) > 1 {
		target = strings.ToLower(os.Args[1])
	}
	if target != "ibge" && target != "ans" && target != "all" {
		log.Fatalf("alvo inválido: %q (use ibge, ans ou all)", os.Args[1])
	}

	cfg := config.LoadConfig()
	if err := cfg.ValidateImport(target); err != nil {
		log.Fatalf("configuração inválida: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Minute,
	)
	defer cancel()

	db, err := postgres.NewDBConnection(ctx, cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("falha ao conectar com o database: %v", err)
	}
	defer db.Close()

	if target == "ibge" || target == "all" {
		importIBGE(ctx, db, cfg)
	}

	if target == "ans" || target == "all" {
		importANS(ctx, db, cfg)
	}
}

func importIBGE(ctx context.Context, db *pgxpool.Pool, cfg *config.Config) {
	ibgeClient := ibge.NewIbgeClient(cfg.BaseUrlIBGE, cfg.BaseUrlLocalidades, nil)

	importStart := time.Now()

	failures := make([]string, 0)

	stateRepository := postgres.NewStateRepo(db)
	stateService := usecases.NewStateUseCase(stateRepository, ibgeClient)

	stageStart := time.Now()
	if err := stateService.Import(ctx); err != nil {
		failures = append(failures, fmt.Sprintf("estados: %v", err))
		log.Printf("falha ao importar estados: %v", err)
	} else {
		log.Printf("etapa estados: %s", time.Since(stageStart))
	}

	cityRepository := postgres.NewCityRepository(db)
	cityService := usecases.NewCityUsecaseImpl(cityRepository, ibgeClient)

	stageStart = time.Now()
	if err := cityService.Import(ctx); err != nil {
		failures = append(failures, fmt.Sprintf("cidades: %v", err))
		log.Printf("falha ao importar cidades: %v", err)
	} else {
		log.Printf("etapa cidades: %s", time.Since(stageStart))
	}

	populationRepository := postgres.NewPopulationRepo(db)
	populationService := usecases.NewPopulationUsecaseImpl(populationRepository, ibgeClient)

	stageStart = time.Now()
	if err := populationService.Import2022(ctx); err != nil {
		failures = append(failures, fmt.Sprintf("população: %v", err))
		log.Printf("falha ao importar população: %v", err)
	} else {
		log.Printf("etapa população: %s", time.Since(stageStart))
	}

	incomeRepository := postgres.NewIncomeRepo(db)
	incomeService := usecases.NewIncomeUsecaseImpl(incomeRepository, ibgeClient)

	stageStart = time.Now()
	if err := incomeService.Import(ctx); err != nil {
		failures = append(failures, fmt.Sprintf("renda: %v", err))
		log.Printf("falha ao importar renda: %v", err)
	} else {
		log.Printf("etapa renda: %s", time.Since(stageStart))
	}

	gdpRepository := postgres.NewGDPRepo(db)
	gdpUseCase := usecases.NewGDPUsecaseImpl(
		gdpRepository,
		ibgeClient,
	)

	stageStart = time.Now()
	if err := gdpUseCase.Import(ctx); err != nil {
		failures = append(failures, fmt.Sprintf("PIB: %v", err))
		log.Printf("falha ao importar PIB: %v", err)
	} else {
		log.Printf("etapa PIB: %s", time.Since(stageStart))
	}

	ageRepository := postgres.NewAgeRepo(db)
	ageUsecase := usecases.NewAgeUsecaseImpl(ageRepository, ibgeClient)

	stageStart = time.Now()
	if err := ageUsecase.Import(ctx); err != nil {
		failures = append(failures, fmt.Sprintf("faixa etária: %v", err))
		log.Printf("falha ao importar faixa etária: %v", err)
	} else {
		log.Printf("etapa faixa etária: %s", time.Since(stageStart))
	}

	if len(failures) > 0 {
		log.Fatalf("importação IBGE incompleta: %s", strings.Join(failures, "; "))
	}

	log.Printf("IBGE import successfully em %s", time.Since(importStart))
}

func importANS(ctx context.Context, db *pgxpool.Pool, cfg *config.Config) {
	ansClient := ans.NewClient(cfg.BaseUrlANS)
	ansRepository := postgres.NewANSRepo(db)
	ansUsecase := usecases.NewANSUsecase(ansClient, ansRepository, cfg.AnsFonte)

	importStart := time.Now()
	if err := ansUsecase.Import(ctx); err != nil {
		log.Fatalf("falha ao importar ANS: %v", err)
	}

	log.Printf("ANS import successfully em %s", time.Since(importStart))
}
