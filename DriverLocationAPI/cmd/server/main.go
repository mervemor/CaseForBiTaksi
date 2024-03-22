package main

import (
	"DriverLocationAPI/internal/handler"
	"DriverLocationAPI/internal/infra/mongodb"
	"DriverLocationAPI/internal/repository"
	"DriverLocationAPI/internal/service"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	ctx := context.Background()
	client, _ := mongodb.ConnectToMongoDB()
	collection := mongodb.GetCollectionFromMongoDB(client)
	mongodb.WriteCSVDataToMongoDB(ctx, "coordinates.csv", collection)

	driverRepo := repository.NewDriverRepository(collection)
	driverService := service.NewDriverService(driverRepo)
	driverHandler := handler.NewDriverHandler(driverService)

	router := mux.NewRouter()

	router.HandleFunc("/find-nearest-driver", driverHandler.DriverHandler).Methods("POST")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("an error occurred when listening and serving HTTP requests.")
	}
}
