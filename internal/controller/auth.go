package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

const prefix = "Bearer "

func auth(JWTsecretKey string, needAdminCheck bool, logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Info("authorization header is missing", "method", r.Method, "path", r.URL.Path)
			http.Error(w, "you're unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, prefix) {
			logger.Info("authorization header has invalid format", "method", r.Method, "path", r.URL.Path)
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		signedToken := strings.TrimPrefix(authHeader, prefix)
		if signedToken == "" {
			logger.Info("authorization token is empty", "method", r.Method, "path", r.URL.Path)
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(signedToken, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(JWTsecretKey), nil
		})

		if err != nil || !token.Valid {
			logger.Info("authorization token is invalid", "method", r.Method, "path", r.URL.Path, "error", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Info("authorization token has invalid claims", "method", r.Method, "path", r.URL.Path)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			logger.Info("authorization token role is missing", "method", r.Method, "path", r.URL.Path)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if needAdminCheck && role != domain.UsersRoleAdmin {
			logger.Info("user does not have admin permissions", "method", r.Method, "path", r.URL.Path)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
