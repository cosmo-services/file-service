package auth

import "errors"

var ErrInvalidCredentials = errors.New("INVALID_CREDENTIALS")

type JwtClient interface {
	ValidateToken(tokenStr string) (*JwtPayload, error)
}

type JwtPayload struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
