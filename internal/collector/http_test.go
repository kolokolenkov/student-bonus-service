package collector

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"student-bonus-service/internal/bonus"
)

func TestHTTPCollector_Send(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем HTTP method.
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		// Проверяем путь.
		if r.URL.Path != "/web/api/closed/bonus/set-data" {
			t.Errorf("expected path /web/api/closed/bonus/set-data, got %s", r.URL.Path)
		}

		// Проверяем Authorization.
		expectedAuthorization := "Bearer test-token"

		if got := r.Header.Get("Authorization"); got != expectedAuthorization {
			t.Errorf(
				"expected Authorization %q, got %q",
				expectedAuthorization,
				got,
			)
		}

		// Проверяем Content-Type.
		contentType := r.Header.Get("Content-Type")

		if contentType != "application/x-www-form-urlencoded" {
			t.Errorf(
				"expected Content-Type application/x-www-form-urlencoded, got %q",
				contentType,
			)
		}

		// Читаем form-data.
		if err := r.ParseForm(); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}

		// Проверяем period.
		if got := r.FormValue("period"); got != "2026-09" {
			t.Errorf(
				"expected period 2026-09, got %s",
				got,
			)
		}

		// Получаем JSON из поля data.
		dataValue := r.FormValue("data")

		if dataValue == "" {
			t.Fatal("data form field is empty")
		}

		// Декодируем JSON.
		var items []bonusItem

		if err := json.Unmarshal([]byte(dataValue), &items); err != nil {
			t.Fatalf("failed to decode data JSON: %v", err)
		}

		// Проверяем количество записей.
		if len(items) != 2 {
			t.Fatalf(
				"expected 2 bonus items, got %d",
				len(items),
			)
		}

		// Проверяем первую запись.
		if items[0].KNumber != 10001 {
			t.Errorf(
				"expected first k_number 10001, got %d",
				items[0].KNumber,
			)
		}

		if items[0].Summa != 9900 {
			t.Errorf(
				"expected first summa 9900, got %v",
				items[0].Summa,
			)
		}

		if items[0].BonusType != "Стипендия" {
			t.Errorf(
				"expected first bonus_type Стипендия, got %s",
				items[0].BonusType,
			)
		}

		// Проверяем вторую запись.
		if items[1].KNumber != 10001 {
			t.Errorf(
				"expected second k_number 10001, got %d",
				items[1].KNumber,
			)
		}

		if items[1].Summa != 32500 {
			t.Errorf(
				"expected second summa 32500, got %v",
				items[1].Summa,
			)
		}

		if items[1].BonusType != "Оплата за практические занятия" {
			t.Errorf(
				"expected second bonus_type Оплата за практические занятия, got %s",
				items[1].BonusType,
			)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		BaseURL:  server.URL,
		SendPath: "/web/api/closed/bonus/set-data",
		Token:    "test-token",
	}

	collector := NewHTTPCollector(
		config,
		server.Client(),
	)

	period := "2026-09"

	calculations := []bonus.BonusCalculation{
		{
			KNumber:   10001,
			Period:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			BonusType: bonus.BonusScholarship,
			Hours:     55,
			Rate:      180,
			Amount:    9900,
			Status:    bonus.StatusCalculated,
		},
		{
			KNumber:   10001,
			Period:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			BonusType: bonus.BonusPractice,
			Hours:     130,
			Rate:      250,
			Amount:    32500,
			Status:    bonus.StatusCalculated,
		},
	}

	err := collector.Send(
		context.Background(),
		period,
		calculations,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestHTTPCollector_Send_ReturnsErrorOnNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

		_, _ = io.WriteString(w, "collector error")
	}))
	defer server.Close()

	config := Config{
		BaseURL:  server.URL,
		SendPath: "/web/api/closed/bonus/set-data",
		Token:    "test-token",
	}

	collector := NewHTTPCollector(
		config,
		server.Client(),
	)

	calculations := []bonus.BonusCalculation{
		{
			KNumber:   10001,
			Period:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			BonusType: bonus.BonusScholarship,
			Hours:     55,
			Rate:      180,
			Amount:    9900,
			Status:    bonus.StatusCalculated,
		},
	}

	err := collector.Send(
		context.Background(),
		"2026-09",
		calculations,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHTTPCollector_Send_FormEncoding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("failed to parse form body: %v", err)
		}

		if values.Get("period") != "2026-09" {
			t.Errorf(
				"expected period 2026-09, got %s",
				values.Get("period"),
			)
		}

		if values.Get("data") == "" {
			t.Error("expected data field to be present")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		BaseURL:  server.URL,
		SendPath: "/web/api/closed/bonus/set-data",
		Token:    "test-token",
	}

	collector := NewHTTPCollector(
		config,
		server.Client(),
	)

	calculations := []bonus.BonusCalculation{
		{
			KNumber:   10001,
			Period:    time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
			BonusType: bonus.BonusScholarship,
			Hours:     55,
			Rate:      180,
			Amount:    9900,
			Status:    bonus.StatusCalculated,
		},
	}

	err := collector.Send(
		context.Background(),
		"2026-09",
		calculations,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
