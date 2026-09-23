package bonus

import "context"

type Collector interface {
	Send(
		ctx context.Context,
		period string,
		calculations []BonusCalculation,
	) error
}
