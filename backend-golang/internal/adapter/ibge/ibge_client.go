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

const userAgent = "trads-market-test/1.0 (+https://github.com/CunhazadanoDale/trads-market-test)"

var retryBackoffs = []time.Duration{
	500 * time.Millisecond,
	1 * time.Second,
	2 * time.Second,
}

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
	return c.getFrom(ctx, c.baseURL, path, query, target)
}

func (c *Client) GetFromBase(ctx context.Context, path string, query url.Values, target any) error {
	return c.getFrom(ctx, c.baseURLLocalidades, path, query, target)
}

func (c *Client) getFrom(ctx context.Context, baseURL string, path string,
	query url.Values, target any) error {
	endpoint := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")

	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	for attempt := 0; ; attempt++ {
		retryable, err := c.fetch(ctx, endpoint, target)
		if err == nil {
			return nil
		}

		if !retryable || attempt >= len(retryBackoffs) {
			return err
		}

		if waitErr := waitBackoff(ctx, retryBackoffs[attempt]); waitErr != nil {
			return fmt.Errorf("tentativa do IBGE interrompida: %w", waitErr)
		}
	}
}

func (c *Client) fetch(ctx context.Context, endpoint string, target any) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("criar requisição do IBGE: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return true, fmt.Errorf("executar requisição do IBGE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		retryable := resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode >= http.StatusInternalServerError
		return retryable, fmt.Errorf("a API do IBGE retornou o status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return false, fmt.Errorf("decodificar resposta do IBGE: %w", err)
	}

	return false, nil
}

func waitBackoff(ctx context.Context, backoff time.Duration) error {
	timer := time.NewTimer(backoff)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
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
