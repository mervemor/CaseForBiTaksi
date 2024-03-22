package repository

import (
	"DriverLocationAPI/internal/domain"
	"DriverLocationAPI/internal/helpers"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DriverRepository interface {
	FindNearestDriver(ctx context.Context, userRadius float64, userCoordinates []float64) ([]domain.DistanceBetweenDriverAndUser, error)
}

type Driver struct {
	Collection *mongo.Collection
}

func NewDriverRepository(collection *mongo.Collection) *Driver {
	return &Driver{
		Collection: collection,
	}
}

func (d *Driver) FindNearestDriver(ctx context.Context, userRadius float64, userCoordinates []float64) ([]domain.DistanceBetweenDriverAndUser, error) {

	/*var records []domain.DistanceDriversFromUser

	filter, _ := d.Collection.Find(ctx, bson.M{"location": bson.M{
		"$nearSphere": bson.M{
			"$geometry": bson.M{
				"type":        "Point",
				"coordinates": []float64{userCoordinates[0], userCoordinates[1]},
			},
			"$maxDistance": userRadius,
		},
	}})

	cursor, err := d.Collection.Find(ctx, filter)
	if err != nil {
		return []domain.DistanceDriversFromUser{}, err
	}
	defer func() {
		if cursor != nil {
			cursor.Close(ctx)
		}
	}()

	if cursor == nil {
		return []domain.DistanceDriversFromUser{}, errors.New("cursor is nil")
	}

	for cursor.Next(ctx) {
		var driverLocation domain.Driver
		if err := cursor.Decode(&driverLocation); err != nil {
			return []domain.DistanceDriversFromUser{}, err
		}

		distance := helpers.Haversine(driverLocation.Location.Coordinates[0], driverLocation.Location.Coordinates[1], userCoordinates[0], userCoordinates[1])
		if distance <= userRadius {
			distanceData := domain.DistanceDriversFromUser{
				DriverID: driverLocation.Id,
				Distance: distance,
			}
			records = append(records, distanceData)
		}
	}

	if err := cursor.Err(); err != nil {
		return []domain.DistanceDriversFromUser{}, err
	}

	return records, nil*/
	var records []domain.DistanceBetweenDriverAndUser

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
		return []domain.DistanceBetweenDriverAndUser{}, err
	}

	defer func() {
		if cursor != nil {
			cursor.Close(ctx)
		}
	}()

	for cursor.Next(ctx) {
		var driverLocation domain.Driver
		if err := cursor.Decode(&driverLocation); err != nil {
			return []domain.DistanceBetweenDriverAndUser{}, err
		}

		distance := helpers.Haversine(driverLocation.Location.Coordinates[0], driverLocation.Location.Coordinates[1], userCoordinates[0], userCoordinates[1])
		if distance <= userRadius {
			distanceData := domain.DistanceBetweenDriverAndUser{
				DriverID: driverLocation.ID,
				Distance: distance,
			}
			records = append(records, distanceData)
		}
	}

	if err := cursor.Err(); err != nil {
		return []domain.DistanceBetweenDriverAndUser{}, err
	}

	return records, nil
}
