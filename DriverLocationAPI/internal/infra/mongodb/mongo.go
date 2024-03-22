package mongodb

import (
	"DriverLocationAPI/internal/domain"
	"context"
	"encoding/csv"
	"errors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"os"
	"strconv"
)

const (
	FileName       = ".env"
	MongoURI       = "MONGOURI"
	DatabaseName   = "BiTaksiCase"
	CollectionName = "driverLocations"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	err := godotenv.Load(FileName)
	if err != nil {
		return nil, errors.New("error loading .env file")
	}

	mongoURI := os.Getenv(MongoURI)

	clientOption := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetCollectionFromMongoDB(client *mongo.Client) *mongo.Collection {
	collection := client.Database(DatabaseName).Collection(CollectionName)
	return collection
}

func WriteCSVDataToMongoDB(ctx context.Context, csvFilePath string, collection *mongo.Collection) error {
	var collectionIsNull domain.GeoJSONLocation
	collectionErr := collection.FindOne(ctx, bson.M{}).Decode(&collectionIsNull)

	if collectionErr != nil {
		file, err := os.Open(csvFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.FieldsPerRecord = 2

		if _, readerErr := reader.Read(); readerErr != nil {
			return readerErr
		}

		var data []interface{}

		for {
			record, readerErr := reader.Read()
			if readerErr == io.EOF {
				break
			}
			if readerErr != nil {
				return err
			}

			latitude, latParseErr := strconv.ParseFloat(record[0], 64)
			if latParseErr != nil {
				return latParseErr
			}

			longitude, longParseErr := strconv.ParseFloat(record[1], 64)
			if longParseErr != nil {
				return longParseErr
			}

			geoJSONlocation := domain.GeoJSONLocation{
				Type:        "Point",
				Coordinates: []float64{latitude, longitude},
			}

			driver := domain.Driver{
				ID:       primitive.NewObjectID(),
				Location: geoJSONlocation,
			}

			data = append(data, driver)
		}

		_, err = collection.InsertMany(ctx, data)
		if err != nil {
			return err
		}

		indexes := []mongo.IndexModel{
			{
				Keys: bson.M{"location": "2dsphere"},
			},
		}

		_, err = collection.Indexes().CreateMany(ctx, indexes)
		if err != nil {
			return err
		}
	}

	return nil
}
