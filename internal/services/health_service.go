package services

import (
	"context"

	"github.com/alvarolucio2007/Scholarly/internal/ports"
)

type HealthService struct {
	checkers []ports.HealthChecker
}

func NewHealthService(checkers ...ports.HealthChecker) *HealthService {
	return &HealthService{checkers: checkers}
}

func (s *HealthService) Health(ctx context.Context) error {
	for _, c := range s.checkers {
		if err := c.Health(ctx); err != nil {
			return err
		}
	}
	return nil
}
