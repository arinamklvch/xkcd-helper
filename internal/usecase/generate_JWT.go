package usecase

import (
	"errors"
	"time"

	"github.com/arinamklvch/xkcd-helper/internal/adapter"
	"github.com/arinamklvch/xkcd-helper/internal/dto"
	"github.com/golang-jwt/jwt/v4"
)

const secretKey = "SecretKey"
const tokenTTL = 15 * time.Minute

var ErrInvalidCredentials = errors.New("invalid login or password")

func (s *Service) GenerateJWT(loginRequest dto.LoginRequest) (string, error) {
	user, err := s.usersStorage.GetUser(loginRequest.Login, loginRequest.Password)
	if errors.Is(err, adapter.ErrUserNotFound) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(tokenTTL).Unix()
	claims["role"] = user.Role
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
