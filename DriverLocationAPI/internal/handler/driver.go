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
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := r.Header.Get("Apikey")
	authenticated, err := helpers.TokenAuthenticate(apiKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Internal Server Error:", err)
		return
	}

	if !authenticated {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	var requestPayload domain.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	coordinates := requestPayload.UserCoordinates
	radius := requestPayload.UserRadius

	nearestDrivers, err := d.Service.DriverService(radius, coordinates)
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
