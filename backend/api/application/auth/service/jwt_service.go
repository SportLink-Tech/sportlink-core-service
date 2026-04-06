package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	Generate(accountID string) (string, error)
}

type jwtService struct {
	secret string
}

func NewJWTService(secret string) JWTService {
	return &jwtService{secret: secret}
}

func (s *jwtService) Generate(accountID string) (string, error) {
	claims := jwt.MapClaims{
		"account_id": accountID,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, nil
}
