package student

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByKNumber(
	ctx context.Context,
	kNumber int64,
) (*Student, error) {
	const query = `
		SELECT
			id,
			k_number,
			fio,
			period,
			student_contract_number,
			student_contract_date,
			start_date_of_training,
			end_date_of_training,
			employment_contract,
			location_code,
			department_code,
			wage_rate,
			created_at,
			updated_at
		FROM students
		WHERE k_number = $1
	`

	var student Student

	err := r.db.QueryRowContext(
		ctx,
		query,
		kNumber,
	).Scan(
		&student.ID,
		&student.KNumber,
		&student.FIO,
		&student.Period,
		&student.StudentContractNumber,
		&student.StudentContractDate,
		&student.StartDateOfTraining,
		&student.EndDateOfTraining,
		&student.EmploymentContract,
		&student.LocationCode,
		&student.DepartmentCode,
		&student.WageRate,
		&student.CreatedAt,
		&student.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student %d not found", kNumber)
		}

		return nil, fmt.Errorf("get student: %w", err)
	}

	return &student, nil
}
