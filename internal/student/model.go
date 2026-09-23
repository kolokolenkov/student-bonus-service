package student

import "time"

type Student struct {
	ID                    int64
	KNumber               int64
	FIO                   string
	Period                *time.Time
	StudentContractNumber *string
	StudentContractDate   *time.Time
	StartDateOfTraining   time.Time
	EndDateOfTraining     *time.Time
	EmploymentContract    *string
	LocationCode          *string
	DepartmentCode        *string
	WageRate              *float64
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
