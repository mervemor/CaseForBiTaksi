package main

import (
	"MatchingAPI/internal/handler"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/find-driver", handler.RiderHandler).Methods("POST")

	err := http.ListenAndServe(":8081", router)
	if err != nil {
		fmt.Println("an error occurred when listening and serving HTTP requests.")
	}
}
