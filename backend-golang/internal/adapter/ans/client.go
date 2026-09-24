package ans

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
)

const (
	expectedColumns = 16
	minPeriod       = 1900
	downloadTimeout = 15 * time.Minute
)

const DefaultFileURL = "https://dadosabertos.ans.gov.br/FTP/PDA/taxa_de_cobertura_de_planos_de_saude-047/pda-047-taxa_cobertura.csv"

var requiredHeader = []string{"PERIODO", "CD_MUNICIPIO"}

type Client struct {
	DefaultFileURL string
	HTTPClient     *http.Client
}

func NewClient(fileURL string) *Client {
	if fileURL == "" {
		fileURL = DefaultFileURL
	}

	return &Client{
		DefaultFileURL: fileURL,
		HTTPClient:     &http.Client{Timeout: downloadTimeout},
	}
}

func (c *Client) Fetch(ctx context.Context) ([]domain.ANSBeneficiaryRow, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.DefaultFileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ANS: montar requisição do dataset: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ANS: baixar dataset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ANS: HTTP %d ao baixar dataset", resp.StatusCode)
	}

	return c.Parse(resp.Body)
}

func (c *Client) Parse(r io.Reader) ([]domain.ANSBeneficiaryRow, error) {
	reader := csv.NewReader(bufio.NewReader(r))
	reader.Comma = ';'
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("ANS: ler cabeçalho do dataset: %w", err)
	}

	for i, name := range requiredHeader {
		if len(header) <= i || strings.TrimSpace(header[i]) != name {
			return nil, fmt.Errorf("ANS: cabeçalho inesperado na coluna %d (esperado %q)", i+1, name)
		}
	}

	rows := make([]domain.ANSBeneficiaryRow, 0)
	line := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		line++

		if err != nil {
			log.Printf("ANS: linha %d inválida (%v) — pulada", line, err)
			continue
		}

		row, ok := parseLine(record)
		if !ok {
			log.Printf("ANS: linha %d com valores inválidos — pulada", line)
			continue
		}

		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("ANS: nenhuma linha válida encontrada no dataset")
	}

	return rows, nil
}

func parseLine(fields []string) (domain.ANSBeneficiaryRow, bool) {
	if len(fields) < expectedColumns {
		return domain.ANSBeneficiaryRow{}, false
	}

	period, err := strconv.Atoi(normalizeInteger(fields[0]))
	if err != nil || period < minPeriod {
		return domain.ANSBeneficiaryRow{}, false
	}

	code, err := strconv.ParseInt(normalizeInteger(fields[1]), 10, 64)
	if err != nil || code < 0 {
		return domain.ANSBeneficiaryRow{}, false
	}

	beneficiaries, err := strconv.ParseInt(normalizeInteger(fields[9]), 10, 64)
	if err != nil || beneficiaries < 0 {
		return domain.ANSBeneficiaryRow{}, false
	}

	return domain.ANSBeneficiaryRow{
		Year:          period,
		IBGECode:      code,
		Beneficiaries: beneficiaries,
	}, true
}

func normalizeInteger(value string) string {
	value = strings.ReplaceAll(strings.TrimSpace(value), ".", "")
	if idx := strings.Index(value, ","); idx >= 0 {
		value = value[:idx]
	}
	return value
}
