package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecret")

type SigninRequest struct {
	Password string `json:"password"`
}

type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, SigninResponse{Error: "Неверный формат запроса"})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" || req.Password != expectedPassword {
		writeJSON(w, http.StatusUnauthorized, SigninResponse{Error: "Неверный пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passhash": expectedPassword,
		"exp":      time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, SigninResponse{Error: "Ошибка создания токена"})
		return
	}

	writeJSON(w, http.StatusOK, SigninResponse{Token: tokenString})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["passhash"] != pass {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			next(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}
