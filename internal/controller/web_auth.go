package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func webAuth(JWTsecretKey string, logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("JWT_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				logger.Info("auth cookie is missing", "method", r.Method, "path", r.URL.Path)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			logger.Error("failed to read auth cookie", "method", r.Method, "path", r.URL.Path, "error", err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		if err = verifyToken(cookie.Value, JWTsecretKey); err != nil {
			logger.Info("auth cookie token is invalid", "method", r.Method, "path", r.URL.Path, "error", err)
			clearAuthCookie(w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func verifyToken(signedToken string, JWTsecretKey string) error {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(JWTsecretKey), nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("invalid token: %w", err)
	}

	return nil
}

func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "JWT_token",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
}
