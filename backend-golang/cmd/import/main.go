package main

import (
	"context"
	"fmt"
	"log"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
	"github.com/CunhazadanoDale/trads-market-test/internal/config"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()

	client := ibge.NewIbgeClient(cfg.BaseUrlIBGE, cfg.BaseUrlLocalidades, nil)
	states, err := client.GetStates(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Estados encontrados: %d\n", len(states))

	for _, state := range states {
		fmt.Printf(
			"%s - %s (%s)\n",
			state.Sigla,
			state.Nome,
			state.Regiao.Nome,
		)
	}

	fmt.Println()

	cities, err := client.GetCitiesByState(ctx, "PB")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Municípios da PB: %d\n", len(cities))

	for i, city := range cities {
		if i >= 10 {
			break
		}

		fmt.Printf(
			"%d - %s\n",
			city.ID,
			city.Nome,
		)
	}
}