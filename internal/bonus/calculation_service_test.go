package bonus

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCalculationService_SendToCollector_Success(t *testing.T) {
	period := time.Date(
		2026,
		9,
		1,
		0, 0, 0, 0,
		time.UTC,
	)

	calculations := []BonusCalculation{
		{
			KNumber:   10001,
			Period:    period,
			BonusType: BonusScholarship,
			Hours:     55,
			Rate:      180,
			Amount:    9900,
			Status:    StatusCalculated,
		},
		{
			KNumber:   10001,
			Period:    period,
			BonusType: BonusPractice,
			Hours:     130,
			Rate:      250,
			Amount:    32500,
			Status:    StatusCalculated,
		},
	}

	repository := newFakeCalculationRepository(
		calculations,
	)

	collector := &fakeCollector{}

	service := NewCalculationService(
		nil,
		nil,
		repository,
		collector,
	)

	err := service.SendToCollector(
		context.Background(),
		period,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if !collector.called {
		t.Fatal("expected collector to be called")
	}

	if collector.period != "2026-09" {
		t.Fatalf(
			"expected period 2026-09, got %s",
			collector.period,
		)
	}

	if len(collector.calculations) != 2 {
		t.Fatalf(
			"expected 2 calculations, got %d",
			len(collector.calculations),
		)
	}

	for _, calculation := range calculations {
		key := calculationKey(
			calculation.KNumber,
			calculation.Period,
			calculation.BonusType,
		)

		status, exists := repository.statuses[key]

		if !exists {
			t.Fatalf(
				"status not updated for %s",
				key,
			)
		}

		if status != StatusSent {
			t.Fatalf(
				"expected status SENT, got %s",
				status,
			)
		}
	}
}

func TestCalculationService_SendToCollector_Error(
	t *testing.T,
) {
	period := time.Date(
		2026,
		9,
		1,
		0, 0, 0, 0,
		time.UTC,
	)

	calculations := []BonusCalculation{
		{
			KNumber:   10001,
			Period:    period,
			BonusType: BonusScholarship,
			Hours:     55,
			Rate:      180,
			Amount:    9900,
			Status:    StatusCalculated,
		},
	}

	repository := newFakeCalculationRepository(
		calculations,
	)

	expectedError := errors.New(
		"collector unavailable",
	)

	collector := &fakeCollector{
		err: expectedError,
	}

	service := NewCalculationService(
		nil,
		nil,
		repository,
		collector,
	)

	err := service.SendToCollector(
		context.Background(),
		period,
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if !collector.called {
		t.Fatal("expected collector to be called")
	}

	key := calculationKey(
		10001,
		period,
		BonusScholarship,
	)

	status := repository.statuses[key]

	if status != StatusSendError {
		t.Fatalf(
			"expected status SEND_ERROR, got %s",
			status,
		)
	}
}
