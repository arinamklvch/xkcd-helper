package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

const secretKey = "SecretKey"
const prefix = "Bearer "

func authMiddleware(needAdminCheck bool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "you're unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		signedToken := strings.TrimPrefix(authHeader, prefix)
		if signedToken == "" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(signedToken, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if needAdminCheck && role != domain.UsersRoleAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
