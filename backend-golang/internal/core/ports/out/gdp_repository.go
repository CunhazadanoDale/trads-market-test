package out

import "context"

type GDPRepository interface {
	Upsert(ctx context.Context, ibgeCode int64, year int, gdp float64) error
}
