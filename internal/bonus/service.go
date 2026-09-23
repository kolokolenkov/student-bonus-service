package bonus

import (
	"context"
	"fmt"
	"time"

	"student-bonus-service/internal/bonus/rules"
	"student-bonus-service/internal/student"
)

type Service struct {
	repository *Repository

	scholarshipRule *rules.ScholarshipRule
	practiceRule    *rules.PracticeRule
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,

		scholarshipRule: rules.NewScholarshipRule(),
		practiceRule:    rules.NewPracticeRule(),
	}
}

func (s *Service) GetRate(
	ctx context.Context,
	bonusType BonusType,
	date time.Time,
) (*BonusRate, error) {
	return s.repository.GetRate(
		ctx,
		bonusType,
		date,
	)
}

func (s *Service) Calculate(
	ctx context.Context,
	hours student.StudentHours,
	period time.Time,
) ([]BonusCalculation, error) {
	scholarshipRate, err := s.GetRate(
		ctx,
		BonusScholarship,
		period,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get scholarship rate: %w",
			err,
		)
	}

	practiceRate, err := s.GetRate(
		ctx,
		BonusPractice,
		period,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get practice rate: %w",
			err,
		)
	}

	scholarshipAmount := s.scholarshipRule.Calculate(
		hours.TheoryHours,
		scholarshipRate.Rate,
	)

	practiceAmount := s.practiceRule.Calculate(
		hours.PracticeHours,
		practiceRate.Rate,
	)

	return []BonusCalculation{
		{
			KNumber:   hours.KNumber,
			Period:    period,
			BonusType: BonusScholarship,
			Hours:     hours.TheoryHours,
			Rate:      scholarshipRate.Rate,
			Amount:    scholarshipAmount,
			Status:    StatusCalculated,
		},
		{
			KNumber:   hours.KNumber,
			Period:    period,
			BonusType: BonusPractice,
			Hours:     hours.PracticeHours,
			Rate:      practiceRate.Rate,
			Amount:    practiceAmount,
			Status:    StatusCalculated,
		},
	}, nil
}
