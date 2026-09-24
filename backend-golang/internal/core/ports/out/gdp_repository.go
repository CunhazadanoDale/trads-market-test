package out

import "context"

type GDPUpsert struct {
	IBGECode int64
	Year     int
	GDP      float64
}

type GDPRepository interface {
	UpsertMany(ctx context.Context, rows []GDPUpsert) error
}
