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

	if !checkAuthorization(requestToken) {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	riderRequest := domain.RiderRequest{
		ID:              primitive.NewObjectID(),
		Type:            rider.Type,
		UserCoordinates: rider.UserCoordinates,
		UserRadius:      rider.UserRadius,
	}

	riderRequestJSON, err := json.Marshal(riderRequest)
	if err != nil {
		http.Error(w, "Error encoding rider request", http.StatusInternalServerError)
		log.Println("Error encoding rider request:", err)
		return
	}

	headers := map[string]string{
		"Apikey":       r.Header.Get("Apikey"),
		"Content-Type": "application/json",
	}

	resp, err := helpers.SendHTTPRequest("POST", "http://localhost:8080/find-nearest-driver", riderRequestJSON, headers)
	if err != nil {
		http.Error(w, "Error sending request to Driver Location API", http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	var result domain.RiderResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading rider response", http.StatusInternalServerError)
		log.Println("Error reading rider response:", err)
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, "Error decoding rider response", http.StatusBadRequest)
		log.Println("Error decoding rider response:", err)
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

func checkAuthorization(requestToken string) bool {
	authToken := strings.Fields(requestToken)[1]
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<true>"), nil
	})

	if authenticated, ok := claims["authenticated"].(bool); ok && authenticated {
		return true
	}
	return false
}
