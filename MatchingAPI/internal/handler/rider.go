package handler

import (
	"MatchingAPI/internal/domain"
	"MatchingAPI/internal/helpers"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func RiderHandler(w http.ResponseWriter, r *http.Request) {
	var rider domain.RiderRequest

	if err := json.NewDecoder(r.Body).Decode(&rider); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		log.Println("Invalid JSON payload:", err)
		return
	}

	requestToken := r.Header.Get("Authorization")
	if requestToken == "" || requestToken == "Bearer" || len(strings.Fields(requestToken)) == 1 {
		http.Error(w, "JWT token authorization failed", http.StatusUnauthorized)
		return
	}

	authToken := strings.Fields(requestToken)[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<true>"), nil
	})

	if !helpers.CheckAuthorization(claims) {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	riderRequest := domain.RiderRequest{
		ID:              primitive.NewObjectID(),
		Type:            rider.Type,
		UserCoordinates: rider.UserCoordinates,
		UserRadius:      rider.UserRadius,
	}

	riderRequestDataConvertedJson, err := json.Marshal(riderRequest)
	if err != nil {
		http.Error(w, "Error encoding rider request", http.StatusInternalServerError)
		log.Println("Error encoding rider request:", err)
		return
	}

	headers := map[string]string{
		"Apikey":       r.Header.Get("Apikey"),
		"Content-Type": "application/json",
	}

	resp, err := helpers.SendHTTPRequest("POST", "http://localhost:8080/find-nearest-driver", riderRequestDataConvertedJson, headers)
	if err != nil {
		http.Error(w, "Error sending request to find driver", http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	var result domain.RiderResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response from find driver", http.StatusInternalServerError)
		log.Println("Error reading response from find driver:", err)
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, "Error decoding response from find driver", http.StatusBadRequest)
		log.Println("Error decoding response from find driver:", err)
		return
	}

	switch result.Status {
	case http.StatusNotFound:
		http.Error(w, "404 - Not Found", http.StatusNotFound)
		return
	case http.StatusBadRequest:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
