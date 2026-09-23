package bonus

import (
	"context"
	"fmt"
	"time"

	"student-bonus-service/internal/student"
)

type CalculationService struct {
	studentService        *student.Service
	bonusService          *Service
	calculationRepository CalculationRepositoryInterface
	collector             Collector
}

func NewCalculationService(
	studentService *student.Service,
	bonusService *Service,
	calculationRepository CalculationRepositoryInterface,
	collector Collector,
) *CalculationService {
	return &CalculationService{
		studentService:        studentService,
		bonusService:          bonusService,
		calculationRepository: calculationRepository,
		collector:             collector,
	}
}

func (s *CalculationService) CalculateMonth(
	ctx context.Context,
	period time.Time,
) error {
	from := time.Date(
		period.Year(),
		period.Month(),
		1,
		0, 0, 0, 0,
		period.Location(),
	)

	to := from.AddDate(0, 1, 0)

	hours, err := s.studentService.GetHoursByPeriod(
		ctx,
		from,
		to,
	)
	if err != nil {
		return fmt.Errorf(
			"get student hours: %w",
			err,
		)
	}

	for _, studentHours := range hours {
		calculations, err := s.bonusService.Calculate(
			ctx,
			studentHours,
			from,
		)
		if err != nil {
			return fmt.Errorf(
				"calculate bonus for student %d: %w",
				studentHours.KNumber,
				err,
			)
		}

		for _, calculation := range calculations {
			if err := s.calculationRepository.Create(
				ctx,
				calculation,
			); err != nil {
				return fmt.Errorf(
					"save calculation for student %d: %w",
					studentHours.KNumber,
					err,
				)
			}
		}
	}

	return nil
}

func (s *CalculationService) GetCalculations(
	ctx context.Context,
	period *time.Time,
	kNumber *int64,
	status *BonusStatus,
) ([]BonusCalculation, error) {
	return s.calculationRepository.GetAll(
		ctx,
		period,
		kNumber,
		status,
	)
}

func (s *CalculationService) SendToCollector(
	ctx context.Context,
	period time.Time,
) error {
	calculations, err := s.calculationRepository.GetByPeriod(
		ctx,
		period,
	)
	if err != nil {
		return fmt.Errorf(
			"get calculations: %w",
			err,
		)
	}

	if len(calculations) == 0 {
		return fmt.Errorf(
			"no calculations for period %s",
			period.Format("2006-01"),
		)
	}

	collectorPeriod := period.Format("2006-01")

	if err := s.collector.Send(
		ctx,
		collectorPeriod,
		calculations,
	); err != nil {
		for _, calculation := range calculations {
			_ = s.calculationRepository.UpdateStatus(
				ctx,
				calculation.KNumber,
				calculation.Period,
				calculation.BonusType,
				StatusSendError,
			)
		}

		return fmt.Errorf(
			"send calculations to collector: %w",
			err,
		)
	}

	for _, calculation := range calculations {
		if err := s.calculationRepository.UpdateStatus(
			ctx,
			calculation.KNumber,
			calculation.Period,
			calculation.BonusType,
			StatusSent,
		); err != nil {
			return fmt.Errorf(
				"update calculation status: %w",
				err,
			)
		}
	}

	return nil
}
