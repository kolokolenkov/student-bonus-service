package student

import (
	"context"
	"time"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByKNumber(
	ctx context.Context,
	kNumber int64,
) (*Student, error) {
	return s.repository.GetByKNumber(ctx, kNumber)
}

func (s *Service) GetHoursByPeriod(
	ctx context.Context,
	from time.Time,
	to time.Time,
) ([]StudentHours, error) {
	return s.repository.GetHoursByPeriod(
		ctx,
		from,
		to,
	)
}
