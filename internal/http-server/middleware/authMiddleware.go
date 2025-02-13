package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "middleware.authMiddleware"
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			render.HTML(w, r, fmt.Sprintf("401 Unauthorized"))
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.HTML(w, r, fmt.Sprintf("401 Unauthorized"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, err := extractUserIDFromClaims(claims)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				render.HTML(w, r, fmt.Sprintf("401 Unauthorized"))
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func extractUserIDFromClaims(claims jwt.MapClaims) (int, error) {
	// Проверяем, что поле userID существует
	userIDValue, ok := claims["userID"]
	if !ok {
		return 0, errors.New("userID not found in token")
	}

	// Проверяем тип userID
	switch v := userIDValue.(type) {
	case float64:
		return int(v), nil // JWT числа всегда float64
	case int:
		return v, nil
	default:
		return 0, errors.New("userID is not a number")
	}
}
