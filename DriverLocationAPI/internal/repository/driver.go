package repository

import (
	"DriverLocationAPI/internal/domain"
	"DriverLocationAPI/internal/helpers"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DriverRepository interface {
	FindNearestDriver(ctx context.Context, userRadius float64, userCoordinates []float64) ([]domain.NearestDriver, error)
}

type Driver struct {
	Collection *mongo.Collection
}

func NewDriverRepository(collection *mongo.Collection) *Driver {
	return &Driver{
		Collection: collection,
	}
}

func (d *Driver) FindNearestDriver(ctx context.Context, userRadius float64, userCoordinates []float64) ([]domain.NearestDriver, error) {
	var records []domain.NearestDriver

	cursor, err := d.Collection.Find(ctx, bson.M{"location": bson.M{
		"$nearSphere": bson.M{
			"$geometry": bson.M{
				"type":        "Point",
				"coordinates": userCoordinates,
			},
			"$maxDistance": userRadius,
		},
	}})

	if err != nil {
		return []domain.NearestDriver{}, err
	}

	defer func() {
		if cursor != nil {
			cursor.Close(ctx)
		}
	}()

	for cursor.Next(ctx) {
		var driver domain.Driver
		if err := cursor.Decode(&driver); err != nil {
			return []domain.NearestDriver{}, err
		}

		distance := helpers.Haversine(driver.Location.Coordinates[0], driver.Location.Coordinates[1], userCoordinates[0], userCoordinates[1])
		if distance <= userRadius {
			distanceData := domain.NearestDriver{
				DriverID: driver.ID,
				Distance: distance,
			}
			records = append(records, distanceData)
		}
	}

	if err := cursor.Err(); err != nil {
		return []domain.NearestDriver{}, err
	}

	return records, nil
}
