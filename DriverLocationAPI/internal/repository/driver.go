package repository

import (
	"DriverLocationAPI/internal/domain"
	"DriverLocationAPI/internal/helpers"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DriverRepository interface {
	FindNearestDriver(ctx context.Context, userRadius float64, userCoordinates []float64) ([]domain.NearestDriver, error)
	UpsertDriver(ctx context.Context, drivers []domain.DriverUpsertRequest) error
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
				"coordinates": []float64{userCoordinates[0], userCoordinates[1]},
			},
			"$maxDistance": userRadius * 1000,
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

func (d *Driver) UpsertDriver(ctx context.Context, drivers []domain.DriverUpsertRequest) error {
	for _, driver := range drivers {
		objID, err := primitive.ObjectIDFromHex(driver.Id)
		if err != nil {
			objID = primitive.NewObjectID()
		}
		newDriver := domain.Driver{
			ID:       objID,
			Location: driver.Location,
		}

		existingDriver := d.Collection.FindOne(ctx, bson.M{"_id": newDriver.ID})
		if existingDriver.Err() == nil {
			updateResult, err := d.Collection.ReplaceOne(ctx, bson.M{"_id": newDriver.ID}, newDriver)
			if err != nil {
				return err
			}
			fmt.Println("Updated:", updateResult)
		} else {
			_, err := d.Collection.InsertOne(ctx, newDriver)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
