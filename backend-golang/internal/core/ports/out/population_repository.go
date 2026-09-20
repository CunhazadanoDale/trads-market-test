package out

import "context"

type PopulationRepository interface {
	Upsert(ctx context.Context, ibgeCode int64, year int, value int64, source string) error
}