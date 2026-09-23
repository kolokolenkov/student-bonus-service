package bonus

import (
	"context"
)

type fakeCollector struct {
	called       bool
	period       string
	calculations []BonusCalculation
	err          error
}

func (c *fakeCollector) Send(
	ctx context.Context,
	period string,
	calculations []BonusCalculation,
) error {
	c.called = true
	c.period = period
	c.calculations = calculations

	return c.err
}
