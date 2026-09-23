package student

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) GetHoursByPeriod(
	ctx context.Context,
	from time.Time,
	to time.Time,
) ([]StudentHours, error) {
	const query = `
		SELECT
			k_number,
			SUM(hours) AS practice_hours,
			SUM(initial_training) AS theory_hours
		FROM student_hours
		WHERE period_date >= $1
		  AND period_date < $2
		GROUP BY k_number
		ORDER BY k_number
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		from,
		to,
	)
	if err != nil {
		return nil, fmt.Errorf("get student hours: %w", err)
	}
	defer rows.Close()

	var result []StudentHours

	for rows.Next() {
		var hours StudentHours

		err := rows.Scan(
			&hours.KNumber,
			&hours.PracticeHours,
			&hours.TheoryHours,
		)
		if err != nil {
			return nil, fmt.Errorf("scan student hours: %w", err)
		}

		result = append(result, hours)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student hours: %w", err)
	}

	return result, nil
}
