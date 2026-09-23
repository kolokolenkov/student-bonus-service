package bonus

import (
	"context"
	"fmt"
	"time"
)

type fakeCalculationRepository struct {
	calculations []BonusCalculation
	statuses     map[string]BonusStatus
}

func newFakeCalculationRepository(
	calculations []BonusCalculation,
) *fakeCalculationRepository {
	return &fakeCalculationRepository{
		calculations: calculations,
		statuses:     make(map[string]BonusStatus),
	}
}

func (r *fakeCalculationRepository) GetByPeriod(
	ctx context.Context,
	period time.Time,
) ([]BonusCalculation, error) {
	return r.calculations, nil
}

func (r *fakeCalculationRepository) UpdateStatus(
	ctx context.Context,
	kNumber int64,
	period time.Time,
	bonusType BonusType,
	status BonusStatus,
) error {
	key := calculationKey(
		kNumber,
		period,
		bonusType,
	)

	r.statuses[key] = status

	return nil
}

func calculationKey(
	kNumber int64,
	period time.Time,
	bonusType BonusType,
) string {
	return fmt.Sprintf(
		"%d:%s:%s",
		kNumber,
		period.Format("2006-01-02"),
		bonusType,
	)
}

func (r *fakeCalculationRepository) Create(
	ctx context.Context,
	calculation BonusCalculation,
) error {
	r.calculations = append(r.calculations, calculation)
	return nil
}

func (r *fakeCalculationRepository) GetAll(
	ctx context.Context,
	period *time.Time,
	kNumber *int64,
	status *BonusStatus,
) ([]BonusCalculation, error) {
	return r.calculations, nil
}
