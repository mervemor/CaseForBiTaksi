package helpers

import "github.com/dgrijalva/jwt-go"

func CheckAuthorization(claims jwt.MapClaims) bool {
	if authenticated, ok := claims["authenticated"].(bool); ok && authenticated {
		return true
	}
	return false
}
