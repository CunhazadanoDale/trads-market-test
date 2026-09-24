package out

import "context"

type ANSUpsert struct {
	IBGECode      int64
	Year          int
	Beneficiaries int64
	Source        string
}

type ANSRepository interface {
	UpsertMany(ctx context.Context, rows []ANSUpsert) ([]int64, error)
}
