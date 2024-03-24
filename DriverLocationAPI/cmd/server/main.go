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
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	client, _ := mongodb.ConnectToMongoDB()
	collection := mongodb.GetCollectionFromMongoDB(client)
	err := mongodb.ImportCSVDataToMongoDB(ctx, "coordinates.csv", collection)
	if err != nil {
		return
	}

	driverRepo := repository.NewDriverRepository(collection)
	driverService := service.NewDriverService(driverRepo)
	driverHandler := handler.NewDriverHandler(driverService)

	router := mux.NewRouter()

	router.HandleFunc("/find-nearest-driver", driverHandler.DriverHandler).Methods("POST")
	router.HandleFunc("/upsert-driver", driverHandler.UpsertDriverHandler).Methods("POST")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown

		wait := time.Second * 30

		ctx, cancel := context.WithTimeout(ctx, wait)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Println("an error occurred while shutting down the server:", err)
		}

		os.Exit(0)
	}()

	fmt.Println("HTTP server is listening on port 8080...")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("an error occurred when listening and serving HTTP requests:", err)
	}
}
