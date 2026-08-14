package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func webAuth(JWTsecretKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("JWT_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		if err = verifyToken(cookie.Value, JWTsecretKey); err != nil {
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
		return errors.New("invalid token")
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
