package postgres

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	integrationStateAlfaCode = 999901
	integrationStateBetaCode = 999902
	integrationRegion        = "Integracao"
)

type integrationFixture struct {
	pool    *pgxpool.Pool
	alfaID  int64
	betaID  int64
	popYear int
}

func newIntegrationFixture(t *testing.T) *integrationFixture {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL não definida; teste de integração ignorado")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := NewDBConnection(ctx, databaseURL)
	if err != nil {
		t.Fatalf("conectar no banco de integração: %v", err)
	}

	t.Cleanup(pool.Close)

	fixture := &integrationFixture{pool: pool}
	fixture.clean(t)
	t.Cleanup(func() { fixture.clean(t) })
	fixture.seed(t)

	return fixture
}

func (f *integrationFixture) clean(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	if _, err := f.pool.Exec(ctx,
		`DELETE FROM cities WHERE ibge_code IN (99990001, 99990002, 99990003, 99990004)`,
	); err != nil {
		t.Fatalf("limpar cidades de integração: %v", err)
	}

	if _, err := f.pool.Exec(ctx,
		`DELETE FROM states WHERE ibge_code IN ($1, $2)`,
		integrationStateAlfaCode, integrationStateBetaCode,
	); err != nil {
		t.Fatalf("limpar estados de integração: %v", err)
	}
}

func (f *integrationFixture) seed(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	if err := f.pool.QueryRow(ctx,
		`INSERT INTO states (ibge_code, name, uf, region) VALUES ($1, 'Estado Integracao Alfa', 'Z1', $2) RETURNING id`,
		integrationStateAlfaCode, integrationRegion,
	).Scan(&f.alfaID); err != nil {
		t.Fatalf("criar estado alfa: %v", err)
	}

	if err := f.pool.QueryRow(ctx,
		`INSERT INTO states (ibge_code, name, uf, region) VALUES ($1, 'Estado Integracao Beta', 'Z2', $2) RETURNING id`,
		integrationStateBetaCode, integrationRegion,
	).Scan(&f.betaID); err != nil {
		t.Fatalf("criar estado beta: %v", err)
	}

	if err := f.pool.QueryRow(ctx,
		`SELECT COALESCE(MAX(year), 2022) FROM population_indicators`,
	).Scan(&f.popYear); err != nil {
		t.Fatalf("consultar ano máximo de população: %v", err)
	}

	type citySeed struct {
		code    int64
		name    string
		stateID int64
	}

	cities := []citySeed{
		{99990001, "Cidade Integracao Dois", f.alfaID},
		{99990002, "Cidade Integracao Tres", f.alfaID},
		{99990003, "Cidade Integracao Um", f.alfaID},
		{99990004, "Cidade Integracao Quatro", f.betaID},
	}

	cityIDs := make(map[string]int64, len(cities))

	for _, city := range cities {
		var id int64

		if err := f.pool.QueryRow(ctx,
			`INSERT INTO cities (ibge_code, state_id, name) VALUES ($1, $2, $3) RETURNING id`,
			city.code, city.stateID, city.name,
		).Scan(&id); err != nil {
			t.Fatalf("criar cidade %s: %v", city.name, err)
		}

		cityIDs[city.name] = id
	}

	populations := map[string]int64{
		"Cidade Integracao Dois": 5000,
		"Cidade Integracao Tres": 20000,
		"Cidade Integracao Um":   10000,
	}

	for name, value := range populations {
		if _, err := f.pool.Exec(ctx,
			`INSERT INTO population_indicators (city_id, year, value, source) VALUES ($1, $2, $3, 'teste')`,
			cityIDs[name], f.popYear, value,
		); err != nil {
			t.Fatalf("criar população de %s: %v", name, err)
		}
	}
}

func cityNames(cities []domain.CityWithIndicators) []string {
	names := make([]string, 0, len(cities))

	for _, city := range cities {
		names = append(names, city.City.Name)
	}

	return names
}

func TestIntegrationFindByState(t *testing.T) {
	fixture := newIntegrationFixture(t)

	repo := NewCityRepository(fixture.pool)
	ctx := context.Background()

	t.Run("filtro por nome", func(t *testing.T) {
		cities, total, err := repo.FindByState(ctx, integrationStateAlfaCode, domain.PaginacaoFilter{
			Page: 1, Size: 10, Nome: "Integracao Tres",
		})
		if err != nil {
			t.Fatalf("FindByState: %v", err)
		}

		if total != 1 || len(cities) != 1 {
			t.Fatalf("esperado 1 cidade filtrada, obtido total=%d len=%d", total, len(cities))
		}

		if cities[0].City.Name != "Cidade Integracao Tres" {
			t.Fatalf("cidade inesperada no filtro: %s", cities[0].City.Name)
		}
	})

	t.Run("filtro sem correspondencia", func(t *testing.T) {
		cities, total, err := repo.FindByState(ctx, integrationStateAlfaCode, domain.PaginacaoFilter{
			Page: 1, Size: 10, Nome: "Nao Existe",
		})
		if err != nil {
			t.Fatalf("FindByState: %v", err)
		}

		if total != 0 || len(cities) != 0 {
			t.Fatalf("esperado resultado vazio, obtido total=%d len=%d", total, len(cities))
		}
	})

	t.Run("paginacao", func(t *testing.T) {
		first, totalFirst, err := repo.FindByState(ctx, integrationStateAlfaCode, domain.PaginacaoFilter{
			Page: 1, Size: 2,
		})
		if err != nil {
			t.Fatalf("FindByState página 1: %v", err)
		}

		second, totalSecond, err := repo.FindByState(ctx, integrationStateAlfaCode, domain.PaginacaoFilter{
			Page: 2, Size: 2,
		})
		if err != nil {
			t.Fatalf("FindByState página 2: %v", err)
		}

		if totalFirst != 3 || totalSecond != 3 {
			t.Fatalf("esperado total 3 nas duas páginas, obtido %d e %d", totalFirst, totalSecond)
		}

		if len(first) != 2 || len(second) != 1 {
			t.Fatalf("esperado 2 e 1 cidades por página, obtido %d e %d", len(first), len(second))
		}

		expectedFirst := []string{"Cidade Integracao Dois", "Cidade Integracao Tres"}
		if !slices.Equal(cityNames(first), expectedFirst) {
			t.Fatalf("página 1 inesperada: %v", cityNames(first))
		}

		if second[0].City.Name != "Cidade Integracao Um" {
			t.Fatalf("página 2 inesperada: %v", cityNames(second))
		}
	})

	t.Run("ordenacao por nome", func(t *testing.T) {
		asc, _, err := repo.FindByState(ctx, integrationStateAlfaCode, domain.PaginacaoFilter{
			Page: 1, Size: 10, Ordenar: "", Ordem: "asc",
		})
		if err != nil {
			t.Fatalf("FindByState asc: %v", err)
		}

		expectedAsc := []string{
			"Cidade Integracao Dois",
			"Cidade Integracao Tres",
			"Cidade Integracao Um",
		}
		if !slices.Equal(cityNames(asc), expectedAsc) {
			t.Fatalf("ordenação asc inesperada: %v", cityNames(asc))
		}

		desc, _, err := repo.FindByState(ctx, integrationStateAlfaCode, domain.PaginacaoFilter{
			Page: 1, Size: 10, Ordenar: "", Ordem: "desc",
		})
		if err != nil {
			t.Fatalf("FindByState desc: %v", err)
		}

		expectedDesc := []string{
			"Cidade Integracao Um",
			"Cidade Integracao Tres",
			"Cidade Integracao Dois",
		}
		if !slices.Equal(cityNames(desc), expectedDesc) {
			t.Fatalf("ordenação desc inesperada: %v", cityNames(desc))
		}
	})

	t.Run("ordenacao por populacao", func(t *testing.T) {
		desc, _, err := repo.FindByState(ctx, integrationStateAlfaCode, domain.PaginacaoFilter{
			Page: 1, Size: 10, Ordenar: "populacao", Ordem: "desc",
		})
		if err != nil {
			t.Fatalf("FindByState populacao desc: %v", err)
		}

		expectedDesc := []string{
			"Cidade Integracao Tres",
			"Cidade Integracao Um",
			"Cidade Integracao Dois",
		}
		if !slices.Equal(cityNames(desc), expectedDesc) {
			t.Fatalf("ordenação por população desc inesperada: %v", cityNames(desc))
		}

		asc, _, err := repo.FindByState(ctx, integrationStateAlfaCode, domain.PaginacaoFilter{
			Page: 1, Size: 10, Ordenar: "populacao", Ordem: "asc",
		})
		if err != nil {
			t.Fatalf("FindByState populacao asc: %v", err)
		}

		expectedAsc := []string{
			"Cidade Integracao Dois",
			"Cidade Integracao Um",
			"Cidade Integracao Tres",
		}
		if !slices.Equal(cityNames(asc), expectedAsc) {
			t.Fatalf("ordenação por população asc inesperada: %v", cityNames(asc))
		}
	})

	t.Run("indicadores carregados", func(t *testing.T) {
		cities, _, err := repo.FindByState(ctx, integrationStateAlfaCode, domain.PaginacaoFilter{
			Page: 1, Size: 10,
		})
		if err != nil {
			t.Fatalf("FindByState: %v", err)
		}

		for _, city := range cities {
			if city.Population == nil {
				t.Fatalf("cidade sem população carregada: %s", city.City.Name)
			}
		}
	})
}

func TestIntegrationFindStates(t *testing.T) {
	fixture := newIntegrationFixture(t)

	repo := NewMetricsRepository(fixture.pool)
	ctx := context.Background()

	t.Run("filtra por regiao", func(t *testing.T) {
		states, err := repo.FindStates(ctx, integrationRegion)
		if err != nil {
			t.Fatalf("FindStates: %v", err)
		}

		if len(states) != 2 {
			t.Fatalf("esperado 2 estados na região, obtido %d", len(states))
		}

		if states[0].State.Name != "Estado Integracao Alfa" ||
			states[1].State.Name != "Estado Integracao Beta" {
			t.Fatalf("ordenação de estados inesperada: %s, %s",
				states[0].State.Name, states[1].State.Name)
		}

		if states[0].Municipalities != 3 {
			t.Fatalf("esperado 3 municípios em Alfa, obtido %d", states[0].Municipalities)
		}

		if states[1].Municipalities != 1 {
			t.Fatalf("esperado 1 município em Beta, obtido %d", states[1].Municipalities)
		}

		if states[0].Population == nil {
			t.Fatal("Alfa deveria ter população agregada")
		}

		if states[0].Population.Value != 35000 {
			t.Fatalf("esperado população 35000 em Alfa, obtido %d", states[0].Population.Value)
		}

		if states[1].Population != nil {
			t.Fatalf("Beta não deveria ter população agregada, obtido %d", states[1].Population.Value)
		}
	})

	t.Run("sem regiao retorna todos", func(t *testing.T) {
		states, err := repo.FindStates(ctx, "")
		if err != nil {
			t.Fatalf("FindStates: %v", err)
		}

		if len(states) < 2 {
			t.Fatalf("esperado ao menos 2 estados, obtido %d", len(states))
		}

		codes := make([]int64, 0, len(states))
		for _, state := range states {
			codes = append(codes, state.State.IBGECode)
		}

		if !slices.Contains(codes, int64(integrationStateAlfaCode)) ||
			!slices.Contains(codes, int64(integrationStateBetaCode)) {
			t.Fatal("estados de integração ausentes na listagem sem filtro")
		}
	})

	t.Run("regiao inexistente", func(t *testing.T) {
		states, err := repo.FindStates(ctx, "Nao Existe")
		if err != nil {
			t.Fatalf("FindStates: %v", err)
		}

		if len(states) != 0 {
			t.Fatalf("esperado nenhum estado, obtido %d", len(states))
		}
	})
}
