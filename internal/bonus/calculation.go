package bonus

import "time"

type BonusStatus string

const (
	StatusCalculated BonusStatus = "CALCULATED"
	StatusSent       BonusStatus = "SENT"
	StatusSendError  BonusStatus = "SEND_ERROR"
)

type BonusCalculation struct {
	KNumber int64
	Period  time.Time

	BonusType BonusType

	Hours  float64
	Rate   float64
	Amount float64

	Status BonusStatus
}
