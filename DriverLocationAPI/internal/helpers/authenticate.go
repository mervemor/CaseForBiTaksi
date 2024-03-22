package helpers

import "github.com/dgrijalva/jwt-go"

func TokenAuthenticate(apiKey string) (bool, error) {
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(apiKey, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("Apikey"), nil
	})

	authenticated, ok := claims["authenticated"].(bool)
	if !ok || !authenticated {
		return false, nil
	}

	return true, nil
}
