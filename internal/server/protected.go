package server

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func AdminProtected(w http.ResponseWriter, r *http.Request) (float64, error) {
	userType, userId,  err := ProtectedHandler(w, r)

	if err != nil {
		return  0, err
	}

	if userType != "admin" {
		return  0, NewAPIError(http.StatusForbidden, fmt.Errorf("not authoirized"))
	}

	return userId, nil
}

func ApplicantProtected(w http.ResponseWriter, r *http.Request) ( float64, error) {
	userType, userId,  err := ProtectedHandler(w, r)

	if err != nil {
		return 0, err
	}

	if userType != "applicant" {
		return 0, NewAPIError(http.StatusForbidden, fmt.Errorf("not authoirized"))
	}

	return  userId, nil
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) (string, float64, error) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", 0,  NewAPIError(http.StatusUnauthorized, fmt.Errorf("no authorization header"))
	}

	claims, err := verifyToken(tokenString)
	if err != nil {
		return "", 0,  NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid token"))
	}

	userType, err := getTypeFromClaims(claims)
	if err != nil {
		return "", 0,  err
	}

	userID, err := getIdFromClaims(claims)
	if err != nil {
		return "", 0, err
	}

	return userType, userID, nil
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

func getIdFromClaims(claims jwt.MapClaims) (float64, error) {
	if id, ok := claims["id"].(float64); ok {
		return id, nil
	}
	return 0, fmt.Errorf("type not found in token claims")
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
