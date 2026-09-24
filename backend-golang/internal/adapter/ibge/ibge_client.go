package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/CunhazadanoDale/trads-market-test/internal/adapter/http/dtos"
)

type Client struct {
	baseURL            string
	baseURLLocalidades string
	httpClient         *http.Client
}

func NewIbgeClient(baseURL string, baseURLLocalidades string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &Client{
		baseURL:            baseURL,
		baseURLLocalidades: baseURLLocalidades,
		httpClient:         httpClient,
	}
}

func (c *Client) Get(ctx context.Context, path string,
	query url.Values, target any) error {
	endpoint := strings.TrimRight(c.baseURL, "/") + "/" + strings.TrimLeft(path, "/")

	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("criar IBGE request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executar IBGE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("IBGE API retornou status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decodificar resposta IBGE: %w", err)
	}

	return nil
}

func (c *Client) GetFromBase(ctx context.Context, path string, query url.Values, target any) error {
	endpoint := strings.TrimRight(c.baseURLLocalidades, "/") + "/" + strings.TrimLeft(path, "/")

	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return fmt.Errorf("create IBGE request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request IBGE API: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf(
			"IBGE API returned status %d",
			resp.StatusCode,
		)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode IBGE response: %w", err)
	}

	return nil
}

func (c *Client) GetPopulation2022(ctx context.Context) ([]dtos.PopulationRecord, error) {
	query := url.Values{}
	query.Set("localidades", "N6[all]")

	var response sidraResponse
	err := c.Get(
		ctx,
		"/agregados/4709/periodos/2022/variaveis/93",
		query,
		&response,
	)
	if err != nil {
		return nil, err
	}

	var records []dtos.PopulationRecord
	for _, topLevel := range response {
		for _, resultado := range topLevel.Resultados {
			records = append(records, resultado.Series...)
		}
	}

	return records, nil
}

type sidraResponse []struct {
	Resultados []struct {
		Series []dtos.PopulationRecord `json:"series"`
	} `json:"resultados"`
}
