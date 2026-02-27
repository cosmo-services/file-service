package pkg

import (
	"errors"
	"strings"
)

var (
	ErrMissingAuthHeader   = errors.New("missing authorization header")
	ErrInvalidAuthHeader   = errors.New("invalid authorization header format")
	ErrUnsupportedAuthType = errors.New("unsupported authorization type")
)

func ParseBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", ErrInvalidAuthHeader
	}

	authType := strings.ToLower(parts[0])
	token := strings.TrimSpace(parts[1])

	if token == "" {
		return "", ErrInvalidAuthHeader
	}

	if authType != "bearer" {
		return "", ErrUnsupportedAuthType
	}

	return token, nil
}
