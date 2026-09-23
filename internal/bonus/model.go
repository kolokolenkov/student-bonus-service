package bonus

import "time"

type BonusType string

const (
	BonusScholarship BonusType = "SCHOLARSHIP"
	BonusPractice    BonusType = "PRACTICE"
)

type BonusRate struct {
	ID        int64
	BonusType BonusType
	Rate      float64
	ValidFrom time.Time
	ValidTo   *time.Time
	CreatedAt time.Time
}
