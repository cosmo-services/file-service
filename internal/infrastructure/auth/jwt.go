package auth

import (
	"errors"
	"time"

	"main/internal/config"
	domain "main/internal/domain/auth"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

type JwtClaims struct {
	Payload *domain.JwtPayload `json:"payload"`
	jwt.RegisteredClaims
}

type JwtClient struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJwtClient(env config.Env) domain.JwtClient {
	return &JwtClient{
		secret: env.JwtSecret,
	}
}

func (s *JwtClient) ValidateToken(tokenStr string) (*domain.JwtPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims.Payload, nil
	}

	return nil, errors.New("invalid token")
}

var Module = fx.Options(
	fx.Provide(NewJwtClient),
)
