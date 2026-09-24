package out

import "context"

type AgeUpsert struct {
	IBGECode   int64
	Year       int
	AgeGroup   string
	Population int64
}

type AgeRepository interface {
	UpsertMany(ctx context.Context, rows []AgeUpsert) error
}
