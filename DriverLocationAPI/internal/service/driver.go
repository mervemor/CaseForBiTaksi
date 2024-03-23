package service

import (
	"DriverLocationAPI/internal/domain"
	"DriverLocationAPI/internal/repository"
	"context"
)

type DriverService struct {
	Repository repository.DriverRepository
}

func NewDriverService(repo repository.DriverRepository) *DriverService {
	return &DriverService{
		Repository: repo,
	}
}

func (d *DriverService) DriverService(userRadius float64, userCoordinates []float64) ([]domain.NearestDriver, error) {
	ctx := context.TODO()
	records, err := d.Repository.FindNearestDriver(ctx, userRadius, userCoordinates)
	if err != nil {
		return []domain.NearestDriver{}, err
	}

	return records, nil
}
