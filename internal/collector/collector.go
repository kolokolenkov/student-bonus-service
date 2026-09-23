package collector

import (
	"context"

	"student-bonus-service/internal/bonus"
)

type Collector interface {
	Send(
		ctx context.Context,
		period string,
		calculations []bonus.BonusCalculation,
	) error
}
