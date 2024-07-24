package server

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func AdminProtected(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, error) {
	userType, claims, err := ProtectedHandler(w, r)

	if err != nil {
		return nil, err
	}

	if userType != "admin" {
		return nil, NewAPIError(http.StatusForbidden, fmt.Errorf("not authoirized"))
	}

	return claims, nil
}

func ApplicantProtected(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, error) {
	userType, claims, err := ProtectedHandler(w, r)

	if err != nil {
		return nil, err
	}

	if userType != "applicant" {
		return nil, NewAPIError(http.StatusForbidden, fmt.Errorf("not authoirized"))
	}

	return claims, nil
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) (string, jwt.MapClaims, error) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", nil, NewAPIError(http.StatusUnauthorized, fmt.Errorf("no authorization header"))
	}

	claims, err := verifyToken(tokenString)
	if err != nil {
		return "", nil, NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid token"))
	}

	userType, err := getTypeFromClaims(claims)
	if err != nil {
		return "", nil, err
	}

	return userType, claims, nil
}

func verifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, err
	}
	return nil, fmt.Errorf("error")
}

func getTypeFromClaims(claims jwt.MapClaims) (string, error) {
	if UserType, ok := claims["type"].(string); ok {
		return UserType, nil
	}
	return "", fmt.Errorf("type not found in token claims")
}

func GetUserIDFromJWT(tokenString string) (float64, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		if id, ok := claims["id"].(float64); ok {
			return id, nil
		}
	}
	return 0, NewAPIError(400, fmt.Errorf("error finding user"))
}
