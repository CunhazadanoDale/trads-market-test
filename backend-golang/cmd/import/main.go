package main

import (
	"context"
	"fmt"
	"log"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/ibge"
)

func main() {
	ctx := context.Background()

	client := ibge.NewIbgeClient(nil)

	data, err := client.GetPopulation2022(ctx)
	if err != nil {
		log.Fatalf("Erro ao obter dados da população: %v", err)
	}

	fmt.Printf("Dados da população em 2022: %+v\n", data)

	for i, item := range data {
		if i >= 5 {
			break
		}

		fmt.Printf("Item %d: %+v\n", i+1, item)
	}
}