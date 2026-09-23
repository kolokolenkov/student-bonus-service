package bonus

import (
	"context"
	"time"
)

type CalculationRepositoryInterface interface {
	Create(
		ctx context.Context,
		calculation BonusCalculation,
	) error

	GetAll(
		ctx context.Context,
		period *time.Time,
		kNumber *int64,
		status *BonusStatus,
	) ([]BonusCalculation, error)

	GetByPeriod(
		ctx context.Context,
		period time.Time,
	) ([]BonusCalculation, error)

	UpdateStatus(
		ctx context.Context,
		kNumber int64,
		period time.Time,
		bonusType BonusType,
		status BonusStatus,
	) error
}
