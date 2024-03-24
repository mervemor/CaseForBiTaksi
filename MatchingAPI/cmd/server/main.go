package main

import (
	"MatchingAPI/internal/handler"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	router := mux.NewRouter()

	router.HandleFunc("/find-driver", handler.RiderHandler).Methods("POST")

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:    ":8081",
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

	fmt.Println("HTTP server is listening on port 8081...")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("an error occurred when listening and serving HTTP requests:", err)
	}
}
