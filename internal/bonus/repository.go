package bonus

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetRate(
	ctx context.Context,
	bonusType BonusType,
	date time.Time,
) (*BonusRate, error) {
	const query = `
		SELECT
			id,
			bonus_type,
			rate,
			valid_from,
			valid_to,
			created_at
		FROM bonus_rates
		WHERE bonus_type = $1
		  AND valid_from <= $2
		  AND (
			  valid_to IS NULL
			  OR valid_to >= $2
		  )
		ORDER BY valid_from DESC
		LIMIT 1
	`

	var rate BonusRate

	err := r.db.QueryRowContext(
		ctx,
		query,
		bonusType,
		date,
	).Scan(
		&rate.ID,
		&rate.BonusType,
		&rate.Rate,
		&rate.ValidFrom,
		&rate.ValidTo,
		&rate.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(
				"bonus rate not found: type=%s, date=%s",
				bonusType,
				date.Format("2006-01-02"),
			)
		}

		return nil, fmt.Errorf("get bonus rate: %w", err)
	}

	return &rate, nil
}
