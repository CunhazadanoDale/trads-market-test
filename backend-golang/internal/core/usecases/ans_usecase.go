package usecases

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/in"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
)

const ansMaxCodesInLog = 20

type ANSFetcher interface {
	Fetch(ctx context.Context) ([]domain.ANSBeneficiaryRow, error)
}

var _ in.ANSUsecase = (*ANSUsecaseImpl)(nil)

type ANSUsecaseImpl struct {
	fetcher ANSFetcher
	repo    out.ANSRepository
	fonte   string
}

func NewANSUsecase(fetcher ANSFetcher, repo out.ANSRepository, fonte string) *ANSUsecaseImpl {
	return &ANSUsecaseImpl{
		fetcher: fetcher,
		repo:    repo,
		fonte:   fonte,
	}
}

func (a *ANSUsecaseImpl) Import(ctx context.Context) error {
	if a.fonte == "" {
		return fmt.Errorf("ANS_FONTE não configurada")
	}

	rows, err := a.fetcher.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("buscar dataset da ANS: %w", err)
	}

	aggregated := make([]out.ANSUpsert, 0, len(rows))
	indexes := make(map[domain.ANSKey]int, len(rows))
	for _, row := range rows {
		key := domain.ANSKey{Year: row.Year, IBGECode: row.IBGECode}
		if idx, ok := indexes[key]; ok {
			aggregated[idx].Beneficiaries += row.Beneficiaries
			continue
		}

		indexes[key] = len(aggregated)
		aggregated = append(aggregated, out.ANSUpsert{
			IBGECode:      row.IBGECode,
			Year:          row.Year,
			Beneficiaries: row.Beneficiaries,
			Source:        a.fonte,
		})
	}

	missing, err := a.repo.UpsertMany(ctx, aggregated)
	if err != nil {
		return fmt.Errorf("gravar beneficiários da ANS: %w", err)
	}

	if len(missing) > 0 {
		log.Printf("ANS: %d códigos sem município correspondente — ignorados: %s",
			len(missing), formatANSCodes(missing))
	}

	imported := len(aggregated) - len(missing)
	if imported == 0 {
		return fmt.Errorf("ANS: nenhum município importado (todos os códigos são desconhecidos)")
	}

	log.Printf("ANS: %d municípios importados com sucesso", imported)
	return nil
}

func formatANSCodes(codes []int64) string {
	truncated := len(codes) > ansMaxCodesInLog
	if truncated {
		codes = codes[:ansMaxCodesInLog]
	}

	parts := make([]string, 0, len(codes)+1)
	for _, code := range codes {
		parts = append(parts, strconv.FormatInt(code, 10))
	}
	if truncated {
		parts = append(parts, "...")
	}

	return strings.Join(parts, ", ")
}
