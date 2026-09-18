package ibge

import (
	"context"
	"testing"
)

func TestAggregates(t *testing.T) {
	client := NewIbgeClient(nil)

	aggregates, err := client.GetAggregates(context.Background())
	if err != nil {
		t.Fatalf("Erro ao obter agregados: %v", err)
	}

	if len(aggregates) == 0 {
		t.Fatal("Nenhum agregado retornado")
	}
}