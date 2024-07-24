package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recruit-sys/internal/models"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var SECRET = []byte("this-is-my-secret-for-jwt")

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) error {

	var userLogin models.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		return InvalidJson()
	}

	if err := loginValidate(&userLogin); err != nil {
		return InvalidReqJsonData(err)
	}

	user, err := s.db.SelectUserWhereMail(userLogin.Email)
	if err != nil {
		return err
	}

	if !CheckPasswordHash(userLogin.Password, user.PasswordHash) {
		return NewAPIError(401, fmt.Errorf("unauthorized"))
	}

	token := GetJWT(user.ID, user.UserType)

	cookie := &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("token : " + token)
	return nil
}

func GetJWT(id int, user_type string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"type": user_type,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString(SECRET)
	return tokenString
}

func loginValidate(u *models.UserLogin) map[string]string {
	errors := make(map[string]string)

	if !isValidEmail(u.Email) {
		errors["email"] = ("invalid email address")
	}

	if len(u.Password) < 3 {
		errors["password"] = ("password must be at least 3 characters long")
	}
	return nil
}
