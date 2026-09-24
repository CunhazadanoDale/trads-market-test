package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/CunhazadanoDale/trads-market-test/internal/core/domain"
	"github.com/CunhazadanoDale/trads-market-test/internal/core/ports/out"
)

type fakeANSFetcher struct {
	rows []domain.ANSBeneficiaryRow
	err  error
}

func (f *fakeANSFetcher) Fetch(context.Context) ([]domain.ANSBeneficiaryRow, error) {
	return f.rows, f.err
}

type fakeANSRepository struct {
	upserted []out.ANSUpsert
	missing  []int64
	err      error
}

func (f *fakeANSRepository) UpsertMany(_ context.Context, rows []out.ANSUpsert) ([]int64, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.upserted = rows
	return f.missing, nil
}

func TestANSUsecaseImportAggregatesByMunicipality(t *testing.T) {
	fetcher := &fakeANSFetcher{rows: []domain.ANSBeneficiaryRow{
		{Year: 2026, IBGECode: 120005, Beneficiaries: 150},
		{Year: 2026, IBGECode: 120005, Beneficiaries: 100},
		{Year: 2026, IBGECode: 980000, Beneficiaries: 10},
	}}
	repo := &fakeANSRepository{missing: []int64{980000}}
	uc := NewANSUsecase(fetcher, repo)

	if err := uc.Import(context.Background()); err != nil {
		t.Fatalf("Import() error = %v", err)
	}

	if len(repo.upserted) != 2 {
		t.Fatalf("len(upserted) = %d, want 2", len(repo.upserted))
	}
	if repo.upserted[0].IBGECode != 120005 || repo.upserted[0].Beneficiaries != 250 {
		t.Errorf("upserted[0] = %+v, want {IBGECode: 120005 Beneficiaries: 250}", repo.upserted[0])
	}
	if repo.upserted[0].Year != 2026 {
		t.Errorf("upserted[0].Year = %d, want 2026", repo.upserted[0].Year)
	}
	if repo.upserted[0].Source == "" {
		t.Error("upserted[0].Source = empty, want fonte preenchida")
	}
}

func TestANSUsecaseImportPropagatesFetchError(t *testing.T) {
	fetcher := &fakeANSFetcher{err: errors.New("falha de rede")}
	uc := NewANSUsecase(fetcher, &fakeANSRepository{})

	if err := uc.Import(context.Background()); err == nil {
		t.Fatal("Import() error = nil, want error")
	}
}

func TestANSUsecaseImportErrorsWhenNothingImportable(t *testing.T) {
	fetcher := &fakeANSFetcher{rows: []domain.ANSBeneficiaryRow{
		{Year: 2026, IBGECode: 980000, Beneficiaries: 10},
	}}
	repo := &fakeANSRepository{missing: []int64{980000}}
	uc := NewANSUsecase(fetcher, repo)

	if err := uc.Import(context.Background()); err == nil {
		t.Fatal("Import() error = nil, want error")
	}
}

func TestANSUsecaseImportErrorsOnEmptyDataset(t *testing.T) {
	uc := NewANSUsecase(&fakeANSFetcher{}, &fakeANSRepository{})

	if err := uc.Import(context.Background()); err == nil {
		t.Fatal("Import() error = nil, want error")
	}
}

func TestANSUsecaseImportPropagatesUpsertError(t *testing.T) {
	fetcher := &fakeANSFetcher{rows: []domain.ANSBeneficiaryRow{
		{Year: 2026, IBGECode: 120005, Beneficiaries: 100},
	}}
	repo := &fakeANSRepository{err: errors.New("erro no banco")}
	uc := NewANSUsecase(fetcher, repo)

	if err := uc.Import(context.Background()); err == nil {
		t.Fatal("Import() error = nil, want error")
	}
}
