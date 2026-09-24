package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ out.CityRepository = (*CityRepo)(nil)

type CityRepo struct {
	db *pgxpool.Pool
}

func NewCityRepository(db *pgxpool.Pool) *CityRepo {
	return &CityRepo{
		db: db,
	}
}

func (c *CityRepo) Upsert(ctx context.Context, city *domain.City, stateIBGECode int64) error {
	query := `INSERT INTO cities (
			ibge_code,
			state_id,
			name
		)
		SELECT
			$1,
			id,
			$2
		FROM states
		WHERE ibge_code = $3
		ON CONFLICT (ibge_code)
		DO UPDATE SET
			state_id = EXCLUDED.state_id,
			name = EXCLUDED.name,
			updated_at = NOW()
	`

	_, err := c.db.Exec(
		ctx, query, city.IBGECode, city.Name, stateIBGECode,
	)

	return err
}

func (c *CityRepo) StateExists(ctx context.Context, stateIBGECode int64) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM states WHERE ibge_code = $1)`

	var exists bool

	if err := c.db.QueryRow(
		ctx,
		query,
		stateIBGECode,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf(
			"check state %d: %w",
			stateIBGECode,
			err,
		)
	}

	return exists, nil
}

func (c *CityRepo) FindByState(
	ctx context.Context,
	stateIBGECode int64,
	filter domain.PaginacaoFilter,
) ([]domain.CityWithIndicators, int, error) {
	join, orderBy := citySortClauses(filter)

	from := `FROM cities c
		INNER JOIN states s ON s.id = c.state_id` + join

	where := `WHERE s.ibge_code = $1`
	args := []any{stateIBGECode}

	if filter.Nome != "" {
		where += ` AND c.name ILIKE $2`
		args = append(args, "%"+escapeLike(filter.Nome)+"%")
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) %s %s`, from, where)

	var total int

	if err := c.db.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf(
			"count cities by state: %w",
			err,
		)
	}

	listQuery := fmt.Sprintf(`
		SELECT
			c.id,
			c.ibge_code,
			c.name,
			c.state_id
		%s
		%s
		%s
		LIMIT $%d
		OFFSET $%d
	`,
		from,
		where,
		orderBy,
		len(args)+1,
		len(args)+2,
	)

	offset := (filter.Page - 1) * filter.Size
	listArgs := append([]any{}, args...)
	listArgs = append(listArgs, filter.Size, offset)

	rows, err := c.db.Query(
		ctx,
		listQuery,
		listArgs...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"query cities by state: %w",
			err,
		)
	}
	defer rows.Close()

	cities := make([]domain.City, 0)

	for rows.Next() {
		var city domain.City

		if err := rows.Scan(
			&city.ID,
			&city.IBGECode,
			&city.Name,
			&city.StateID,
		); err != nil {
			return nil, 0, fmt.Errorf(
				"scan city: %w",
				err,
			)
		}

		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf(
			"iterate cities: %w",
			err,
		)
	}

	if len(cities) == 0 {
		return []domain.CityWithIndicators{}, total, nil
	}

	cityIDs := make([]int64, 0, len(cities))
	for _, city := range cities {
		cityIDs = append(cityIDs, city.ID)
	}

	populations, err := c.findPopulationsByCityIDs(ctx, cityIDs)
	if err != nil {
		return nil, 0, err
	}

	incomes, err := c.findIncomesByCityIDs(ctx, cityIDs)
	if err != nil {
		return nil, 0, err
	}

	gdps, err := c.findGDPsByCityIDs(ctx, cityIDs)
	if err != nil {
		return nil, 0, err
	}

	result := make([]domain.CityWithIndicators, 0, len(cities))

	for _, city := range cities {
		result = append(result, domain.CityWithIndicators{
			City:       city,
			Population: populations[city.ID],
			Income:     incomes[city.ID],
			GDP:        gdps[city.ID],
		})
	}

	return result, total, nil
}

func citySortClauses(filter domain.PaginacaoFilter) (join string, orderBy string) {
	direction := "ASC"
	if filter.Ordem == "desc" {
		direction = "DESC"
	}

	switch filter.Ordenar {
	case "populacao":
		return `
		LEFT JOIN population_indicators p
			ON p.city_id = c.id
			AND p.year = (SELECT MAX(year) FROM population_indicators)`,
			"ORDER BY p.value " + direction + " NULLS LAST, c.name ASC"
	case "renda":
		return `
		LEFT JOIN income_indicators i
			ON i.city_id = c.id
			AND i.year = (SELECT MAX(year) FROM income_indicators)`,
			"ORDER BY i.average_income " + direction + " NULLS LAST, c.name ASC"
	case "pib":
		return `
		LEFT JOIN gdp_indicators g
			ON g.city_id = c.id
			AND g.year = (SELECT MAX(year) FROM gdp_indicators)`,
			"ORDER BY g.gdp " + direction + " NULLS LAST, c.name ASC"
	default:
		return "", "ORDER BY c.name " + direction
	}
}

func escapeLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `%`, `\%`)
	value = strings.ReplaceAll(value, `_`, `\_`)
	return value
}

func (c *CityRepo) findPopulationsByCityIDs(
	ctx context.Context,
	cityIDs []int64,
) (map[int64]*domain.Indicator[int64], error) {
	const query = `
		SELECT DISTINCT ON (city_id) city_id, year, value
		FROM population_indicators
		WHERE city_id = ANY($1)
		ORDER BY city_id, year DESC
	`

	rows, err := c.db.Query(ctx, query, cityIDs)
	if err != nil {
		return nil, fmt.Errorf("query populations by city ids: %w", err)
	}
	defer rows.Close()

	indicators := make(map[int64]*domain.Indicator[int64])

	for rows.Next() {
		var cityID int64
		var year int
		var value int64

		if err := rows.Scan(&cityID, &year, &value); err != nil {
			return nil, fmt.Errorf("scan population by city id: %w", err)
		}

		indicators[cityID] = &domain.Indicator[int64]{Year: year, Value: value}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate populations by city ids: %w", err)
	}

	return indicators, nil
}

func (c *CityRepo) findIncomesByCityIDs(
	ctx context.Context,
	cityIDs []int64,
) (map[int64]*domain.Indicator[float64], error) {
	const query = `
		SELECT DISTINCT ON (city_id) city_id, year, average_income
		FROM income_indicators
		WHERE city_id = ANY($1)
		ORDER BY city_id, year DESC
	`

	rows, err := c.db.Query(ctx, query, cityIDs)
	if err != nil {
		return nil, fmt.Errorf("query incomes by city ids: %w", err)
	}
	defer rows.Close()

	indicators := make(map[int64]*domain.Indicator[float64])

	for rows.Next() {
		var cityID int64
		var year int
		var value float64

		if err := rows.Scan(&cityID, &year, &value); err != nil {
			return nil, fmt.Errorf("scan income by city id: %w", err)
		}

		indicators[cityID] = &domain.Indicator[float64]{Year: year, Value: value}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate incomes by city ids: %w", err)
	}

	return indicators, nil
}

func (c *CityRepo) findGDPsByCityIDs(
	ctx context.Context,
	cityIDs []int64,
) (map[int64]*domain.Indicator[float64], error) {
	const query = `
		SELECT DISTINCT ON (city_id) city_id, year, gdp
		FROM gdp_indicators
		WHERE city_id = ANY($1)
		ORDER BY city_id, year DESC
	`

	rows, err := c.db.Query(ctx, query, cityIDs)
	if err != nil {
		return nil, fmt.Errorf("query gdps by city ids: %w", err)
	}
	defer rows.Close()

	indicators := make(map[int64]*domain.Indicator[float64])

	for rows.Next() {
		var cityID int64
		var year int
		var value float64

		if err := rows.Scan(&cityID, &year, &value); err != nil {
			return nil, fmt.Errorf("scan gdp by city id: %w", err)
		}

		indicators[cityID] = &domain.Indicator[float64]{Year: year, Value: value}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate gdps by city ids: %w", err)
	}

	return indicators, nil
}

func (c *CityRepo) FindDetailByIBGECode(ctx context.Context, ibgeCode int64) (*domain.CityDetail, error) {
	const query = `
		SELECT
			c.id,
			c.ibge_code,
			c.name,
			c.state_id,
			s.id,
			s.ibge_code,
			s.name,
			s.uf,
			s.region
		FROM cities c
		INNER JOIN states s ON s.id = c.state_id
		WHERE c.ibge_code = $1
	`

	var detail domain.CityDetail

	if err := c.db.QueryRow(ctx, query, ibgeCode).Scan(
		&detail.City.ID,
		&detail.City.IBGECode,
		&detail.City.Name,
		&detail.City.StateID,
		&detail.State.ID,
		&detail.State.IBGECode,
		&detail.State.Name,
		&detail.State.UF,
		&detail.State.Region,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCityNotFound
		}
		return nil, fmt.Errorf("query city detail %d: %w", ibgeCode, err)
	}

	population, err := c.findPopulation(ctx, detail.City.ID)
	if err != nil {
		return nil, err
	}
	detail.Population = population

	income, err := c.findIncome(ctx, detail.City.ID)
	if err != nil {
		return nil, err
	}
	detail.Income = income

	gdp, err := c.findGDP(ctx, detail.City.ID)
	if err != nil {
		return nil, err
	}
	detail.GDP = gdp

	return &detail, nil
}

func (c *CityRepo) findPopulation(ctx context.Context, cityID int64) (*domain.Indicator[int64], error) {
	const query = `
		SELECT year, value
		FROM population_indicators
		WHERE city_id = $1
		ORDER BY year DESC
		LIMIT 1
	`

	var year int
	var value int64

	err := c.db.QueryRow(ctx, query, cityID).Scan(&year, &value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query population for city %d: %w", cityID, err)
	}

	return &domain.Indicator[int64]{Year: year, Value: value}, nil
}

func (c *CityRepo) findIncome(ctx context.Context, cityID int64) (*domain.Indicator[float64], error) {
	const query = `
		SELECT year, average_income
		FROM income_indicators
		WHERE city_id = $1
		ORDER BY year DESC
		LIMIT 1
	`

	var year int
	var value float64

	err := c.db.QueryRow(ctx, query, cityID).Scan(&year, &value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query income for city %d: %w", cityID, err)
	}

	return &domain.Indicator[float64]{Year: year, Value: value}, nil
}

func (c *CityRepo) findGDP(ctx context.Context, cityID int64) (*domain.Indicator[float64], error) {
	const query = `
		SELECT year, gdp
		FROM gdp_indicators
		WHERE city_id = $1
		ORDER BY year DESC
		LIMIT 1
	`

	var year int
	var value float64

	err := c.db.QueryRow(ctx, query, cityID).Scan(&year, &value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query GDP for city %d: %w", cityID, err)
	}

	return &domain.Indicator[float64]{Year: year, Value: value}, nil
}
