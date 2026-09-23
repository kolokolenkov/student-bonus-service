package bonus

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	calculationService *CalculationService
}

func NewHandler(
	calculationService *CalculationService,
) *Handler {
	return &Handler{
		calculationService: calculationService,
	}
}

type calculateRequest struct {
	Period string `json:"period"`
}

type calculateResponse struct {
	Status string `json:"status"`
	Period string `json:"period"`
}

type calculationResponse struct {
	KNumber   int64       `json:"k_number"`
	Period    string      `json:"period"`
	BonusType BonusType   `json:"bonus_type"`
	Hours     float64     `json:"hours"`
	Rate      float64     `json:"rate"`
	Amount    float64     `json:"amount"`
	Status    BonusStatus `json:"status"`
}

func (h *Handler) Calculate(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var request calculateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	period, err := time.Parse(
		"2006-01-02",
		request.Period,
	)
	if err != nil {
		http.Error(
			w,
			"invalid period format, expected YYYY-MM-DD",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.calculationService.CalculateMonth(
		r.Context(),
		period,
	); err != nil {
		http.Error(
			w,
			"calculation failed",
			http.StatusInternalServerError,
		)
		return
	}

	response := calculateResponse{
		Status: "calculated",
		Period: period.Format("2006-01-02"),
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetCalculations(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var period *time.Time

	if value := r.URL.Query().Get("period"); value != "" {
		parsed, err := time.Parse(
			"2006-01-02",
			value,
		)
		if err != nil {
			http.Error(
				w,
				"invalid period format, expected YYYY-MM-DD",
				http.StatusBadRequest,
			)
			return
		}

		period = &parsed
	}

	var kNumber *int64

	if value := r.URL.Query().Get("k_number"); value != "" {
		parsed, err := strconv.ParseInt(
			value,
			10,
			64,
		)
		if err != nil {
			http.Error(
				w,
				"invalid k_number",
				http.StatusBadRequest,
			)
			return
		}

		kNumber = &parsed
	}

	var status *BonusStatus

	if value := r.URL.Query().Get("status"); value != "" {
		parsed := BonusStatus(value)
		status = &parsed
	}

	calculations, err := h.calculationService.GetCalculations(
		r.Context(),
		period,
		kNumber,
		status,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get calculations",
			http.StatusInternalServerError,
		)
		return
	}

	response := make(
		[]calculationResponse,
		0,
		len(calculations),
	)

	for _, calculation := range calculations {
		response = append(
			response,
			calculationResponse{
				KNumber:   calculation.KNumber,
				Period:    calculation.Period.Format("2006-01-02"),
				BonusType: calculation.BonusType,
				Hours:     calculation.Hours,
				Rate:      calculation.Rate,
				Amount:    calculation.Amount,
				Status:    calculation.Status,
			},
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) SendToCollector(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	periodValue := r.URL.Query().Get("period")

	if periodValue == "" {
		http.Error(
			w,
			"period is required",
			http.StatusBadRequest,
		)
		return
	}

	period, err := time.Parse(
		"2006-01-02",
		periodValue,
	)
	if err != nil {
		http.Error(
			w,
			"invalid period format, expected YYYY-MM-DD",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.calculationService.SendToCollector(
		r.Context(),
		period,
	); err != nil {
		http.Error(
			w,
			"failed to send bonuses to collector",
			http.StatusInternalServerError,
		)
		return
	}

	response := map[string]string{
		"status": "sent",
		"period": period.Format("2006-01-02"),
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}
