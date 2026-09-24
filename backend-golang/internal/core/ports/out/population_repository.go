package out

import "context"

type PopulationUpsert struct {
	IBGECode int64
	Year     int
	Value    int64
	Source   string
}

type PopulationRepository interface {
	UpsertMany(ctx context.Context, rows []PopulationUpsert) error
}
