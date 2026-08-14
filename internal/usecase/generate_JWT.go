package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/arinamklvch/xkcd-helper/internal/adapter"
	"github.com/arinamklvch/xkcd-helper/internal/dto"
	"github.com/golang-jwt/jwt/v4"
)

var ErrInvalidCredentials = errors.New("invalid login or password")

func (s *Service) GenerateJWT(loginRequest dto.LoginRequest) (string, error) {
	user, err := s.usersStorage.GetUser(loginRequest.Login, loginRequest.Password)
	if errors.Is(err, adapter.ErrUserNotFound) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", fmt.Errorf("get user from database: %w", err)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Duration(s.tokenTTL) * time.Minute).Unix()
	claims["role"] = user.Role
	signedToken, err := token.SignedString([]byte(s.JWTsecretKey))
	if err != nil {
		return "", fmt.Errorf("sign JWT token: %w", err)
	}

	return signedToken, nil
}
