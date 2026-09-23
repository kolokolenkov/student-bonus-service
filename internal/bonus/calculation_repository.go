package bonus

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CalculationRepository struct {
	db *sql.DB
}

var _ CalculationRepositoryInterface = (*CalculationRepository)(nil)

func NewCalculationRepository(db *sql.DB) *CalculationRepository {
	return &CalculationRepository{
		db: db,
	}
}

func (r *CalculationRepository) Create(
	ctx context.Context,
	calculation BonusCalculation,
) error {
	const query = `
		INSERT INTO bonus_calculations (
			k_number,
			period,
			bonus_type,
			hours,
			rate,
			amount,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (
			k_number,
			period,
			bonus_type
		)
		DO UPDATE SET
			hours = EXCLUDED.hours,
			rate = EXCLUDED.rate,
			amount = EXCLUDED.amount,
			status = EXCLUDED.status,
			updated_at = NOW()
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		calculation.KNumber,
		calculation.Period,
		calculation.BonusType,
		calculation.Hours,
		calculation.Rate,
		calculation.Amount,
		calculation.Status,
	)

	if err != nil {
		return fmt.Errorf(
			"create bonus calculation: %w",
			err,
		)
	}

	return nil
}

func (r *CalculationRepository) GetAll(
	ctx context.Context,
	period *time.Time,
	kNumber *int64,
	status *BonusStatus,
) ([]BonusCalculation, error) {
	query := `
		SELECT
			k_number,
			period,
			bonus_type,
			hours,
			rate,
			amount,
			status
		FROM bonus_calculations
		WHERE 1 = 1
	`

	args := make([]interface{}, 0)

	if period != nil {
		query += fmt.Sprintf(
			" AND period = $%d",
			len(args)+1,
		)
		args = append(args, *period)
	}

	if kNumber != nil {
		query += fmt.Sprintf(
			" AND k_number = $%d",
			len(args)+1,
		)
		args = append(args, *kNumber)
	}

	if status != nil {
		query += fmt.Sprintf(
			" AND status = $%d",
			len(args)+1,
		)
		args = append(args, *status)
	}

	query += `
		ORDER BY period DESC, k_number, bonus_type
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get bonus calculations: %w",
			err,
		)
	}
	defer rows.Close()

	result := make([]BonusCalculation, 0)

	for rows.Next() {
		var calculation BonusCalculation

		err := rows.Scan(
			&calculation.KNumber,
			&calculation.Period,
			&calculation.BonusType,
			&calculation.Hours,
			&calculation.Rate,
			&calculation.Amount,
			&calculation.Status,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan bonus calculation: %w",
				err,
			)
		}

		result = append(result, calculation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate bonus calculations: %w",
			err,
		)
	}

	return result, nil
}

func (r *CalculationRepository) GetByPeriod(
	ctx context.Context,
	period time.Time,
) ([]BonusCalculation, error) {
	const query = `
		SELECT
			k_number,
			period,
			bonus_type,
			hours,
			rate,
			amount,
			status
		FROM bonus_calculations
		WHERE period = $1
		ORDER BY k_number, bonus_type
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		period,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get calculations by period: %w",
			err,
		)
	}
	defer rows.Close()

	result := make([]BonusCalculation, 0)

	for rows.Next() {
		var calculation BonusCalculation

		if err := rows.Scan(
			&calculation.KNumber,
			&calculation.Period,
			&calculation.BonusType,
			&calculation.Hours,
			&calculation.Rate,
			&calculation.Amount,
			&calculation.Status,
		); err != nil {
			return nil, fmt.Errorf(
				"scan calculation: %w",
				err,
			)
		}

		result = append(result, calculation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate calculations: %w",
			err,
		)
	}

	return result, nil
}

func (r *CalculationRepository) UpdateStatus(
	ctx context.Context,
	kNumber int64,
	period time.Time,
	bonusType BonusType,
	status BonusStatus,
) error {
	const query = `
		UPDATE bonus_calculations
		SET
			status = $1,
			updated_at = NOW()
		WHERE k_number = $2
		  AND period = $3
		  AND bonus_type = $4
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		status,
		kNumber,
		period,
		bonusType,
	)
	if err != nil {
		return fmt.Errorf(
			"update bonus calculation status: %w",
			err,
		)
	}

	return nil
}
