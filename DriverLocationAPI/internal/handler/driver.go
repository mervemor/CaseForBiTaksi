package handler

import (
	"DriverLocationAPI/internal/domain"
	"DriverLocationAPI/internal/helpers"
	"DriverLocationAPI/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type DriverHandler struct {
	Service *service.DriverService
}

func NewDriverHandler(s *service.DriverService) *DriverHandler {
	return &DriverHandler{
		Service: s,
	}
}

func (d *DriverHandler) DriverHandler(w http.ResponseWriter, r *http.Request) {

	tokenAuthentication, err := helpers.TokenAuthenticate(r.Header.Get("Apikey"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Internal Server Error:", err)
		return
	}

	if !tokenAuthentication {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	var requestPayload domain.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	nearestDrivers, err := d.Service.NearestDriverService(requestPayload.UserRadius, requestPayload.UserCoordinates)
	if err != nil {
		http.Error(w, "Error finding driver", http.StatusInternalServerError)
		log.Println("Error finding driver:", err)
		return
	}

	response := domain.DriverResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    nearestDrivers[1],
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func (d *DriverHandler) UpsertDriverHandler(w http.ResponseWriter, r *http.Request) {
	var drivers []domain.DriverUpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&drivers); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err := d.Service.UpsertDriverService(drivers)
	if err != nil {
		http.Error(w, "Error upserting drivers", http.StatusInternalServerError)
		log.Println("Error upserting drivers:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Drivers upserted successfully"))
}
