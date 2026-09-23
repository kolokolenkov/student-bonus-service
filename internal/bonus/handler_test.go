package bonus

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"student-bonus-service/internal/student"
)

func TestHandler_SendToCollector(t *testing.T) {
	period := time.Date(
		2026,
		9,
		1,
		0,
		0,
		0,
		0,
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

	repository := newFakeCalculationRepository(calculations)
	collector := &fakeCollector{}

	// Для SendToCollector studentService и bonusService
	// фактически не используются.
	calculationService := NewCalculationService(
		(*student.Service)(nil),
		(*Service)(nil),
		repository,
		collector,
	)

	handler := NewHandler(calculationService)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/bonuses/send?period=2026-09-01",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.SendToCollector(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d, body: %s",
			http.StatusOK,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !collector.called {
		t.Fatal("expected collector to be called")
	}

	if collector.period != "2026-09" {
		t.Errorf(
			"expected collector period 2026-09, got %s",
			collector.period,
		)
	}

	if len(collector.calculations) != 2 {
		t.Fatalf(
			"expected 2 calculations, got %d",
			len(collector.calculations),
		)
	}

	t.Logf("response body: %s", recorder.Body.String())

	var response map[string]string

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf(
			"failed to decode response: %v, body: %s",
			err,
			recorder.Body.String(),
		)
	}

	if response["status"] != "sent" {
		t.Errorf(
			"expected status sent, got %q",
			response["status"],
		)
	}

	if response["period"] != "2026-09-01" {
		t.Errorf(
			"expected period 2026-09-01, got %q",
			response["period"],
		)
	}

	for _, calculation := range calculations {
		key := calculationKey(
			calculation.KNumber,
			calculation.Period,
			calculation.BonusType,
		)

		status, ok := repository.statuses[key]

		if !ok {
			t.Errorf("status not found for calculation %s", key)
			continue
		}

		if status != StatusSent {
			t.Errorf(
				"expected status SENT for %s, got %s",
				key,
				status,
			)
		}
	}
}

func TestHandler_SendToCollector_MissingPeriod(t *testing.T) {
	repository := newFakeCalculationRepository(nil)
	collector := &fakeCollector{}

	calculationService := NewCalculationService(
		(*student.Service)(nil),
		(*Service)(nil),
		repository,
		collector,
	)

	handler := NewHandler(calculationService)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/bonuses/send",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.SendToCollector(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if collector.called {
		t.Fatal("collector should not be called")
	}
}

func TestHandler_SendToCollector_InvalidPeriod(t *testing.T) {
	repository := newFakeCalculationRepository(nil)
	collector := &fakeCollector{}

	calculationService := NewCalculationService(
		(*student.Service)(nil),
		(*Service)(nil),
		repository,
		collector,
	)

	handler := NewHandler(calculationService)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/bonuses/send?period=2026-99-99",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.SendToCollector(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if collector.called {
		t.Fatal("collector should not be called")
	}
}

func TestHandler_SendToCollector_CollectorError(t *testing.T) {
	period := time.Date(
		2026,
		9,
		1,
		0,
		0,
		0,
		0,
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

	repository := newFakeCalculationRepository(calculations)

	collectorError := context.DeadlineExceeded

	collector := &fakeCollector{
		err: collectorError,
	}

	calculationService := NewCalculationService(
		(*student.Service)(nil),
		(*Service)(nil),
		repository,
		collector,
	)

	handler := NewHandler(calculationService)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/bonuses/send?period=2026-09-01",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.SendToCollector(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}

	if !collector.called {
		t.Fatal("expected collector to be called")
	}

	key := calculationKey(
		10001,
		period,
		BonusScholarship,
	)

	if repository.statuses[key] != StatusSendError {
		t.Errorf(
			"expected status SEND_ERROR, got %s",
			repository.statuses[key],
		)
	}
}
