package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessTokenClaims is the payload of JWTs issued by JWTService (HS256).
// Standard fields exp/iat are carried by [jwt.RegisteredClaims].
type AccessTokenClaims struct {
	AccountID string `json:"account_id"`
	jwt.RegisteredClaims
}

type JWTService interface {
	Generate(accountID string) (string, error)
	// Parse verifies the signature with the service secret and decodes claims into [AccessTokenClaims].
	Parse(token string) (*AccessTokenClaims, error)
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

func (s *jwtService) Parse(token string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
