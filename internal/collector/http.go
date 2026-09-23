package collector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"student-bonus-service/internal/bonus"
)

type HTTPCollector struct {
	baseURL  string
	sendPath string
	token    string
	client   *http.Client
}

func NewHTTPCollector(
	cfg Config,
	client *http.Client,
) *HTTPCollector {
	return &HTTPCollector{
		baseURL:  strings.TrimRight(cfg.BaseURL, "/"),
		sendPath: strings.Trim(cfg.SendPath, "/"),
		token:    cfg.Token,
		client:   client,
	}
}

type bonusItem struct {
	KNumber   int64   `json:"k_number"`
	Summa     float64 `json:"summa"`
	BonusType string  `json:"bonus_type"`
}

func (c *HTTPCollector) Send(
	ctx context.Context,
	period string,
	calculations []bonus.BonusCalculation,
) error {
	items := make([]bonusItem, 0, len(calculations))

	for _, calculation := range calculations {
		items = append(items, bonusItem{
			KNumber: calculation.KNumber,
			Summa:   calculation.Amount,
			BonusType: mapBonusType(
				calculation.BonusType,
			),
		})
	}

	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf(
			"marshal collector data: %w",
			err,
		)
	}

	form := url.Values{}
	form.Set("period", period)
	form.Set("data", string(data))

	requestURL := c.baseURL + "/" + c.sendPath

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		requestURL,
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf(
			"create collector request: %w",
			err,
		)
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+c.token,
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf(
			"send request to collector: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"collector returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

func mapBonusType(
	bonusType bonus.BonusType,
) string {
	switch bonusType {
	case bonus.BonusScholarship:
		return "Стипендия"

	case bonus.BonusPractice:
		return "Оплата за практические занятия"

	default:
		return string(bonusType)
	}
}

var _ Collector = (*HTTPCollector)(nil)
