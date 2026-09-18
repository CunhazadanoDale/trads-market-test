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

const defaultbaseURL = "https://servicodados.ibge.gov.br/api/v3"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewIbgeClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &Client{
		baseURL:    defaultbaseURL,
		httpClient: httpClient,
	}
}

func (c *Client) Get(ctx context.Context, path string, 
	query url.Values, target any) error {
		endpoint := strings.TrimRight(c.baseURL, "/") + "/" + strings.TrimLeft(path,"/")

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
			return fmt.Errorf("IBGE API retornou status %d", resp.StatusCode,)
		}

		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("decodificar resposta IBGE: %w", err)
		}

		return nil
	}

func (c *Client) GetAggregates(ctx context.Context) ([]dtos.Aggregate, error) {
	var response []dtos.Aggregate
	if err := c.Get(ctx, "/agregados", nil, &response); err != nil {
		return nil, err
	}

	return response, nil
}