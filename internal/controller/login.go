package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/arinamklvch/xkcd-helper/internal/dto"
	"github.com/arinamklvch/xkcd-helper/internal/usecase"
)

// @Summary Login
// @Description Validates user credentials and returns a JWT token.
// @Tags auth
// @Accept json
// @Produce plain
// @Param request body dto.LoginRequest true "User login and password"
// @Success 200 {string} string "JWT token"
// @Failure 400 {string} string "Invalid JSON data or validation error"
// @Failure 401 {string} string "Invalid login or password"
// @Failure 500 {string} string "Failed to generate token"
// @Router /login [post]
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	// decode request
	var loginRequest dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			fmt.Println("failed to close response body:", err)
		}
	}()

	// generate JWT token
	signedToken, err := h.service.GenerateJWT(loginRequest)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// set token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "JWT_token",
		Value: signedToken,
	})

	// redirect to search page
	http.Redirect(w, r, "/search-comics", http.StatusMovedPermanently)
}
