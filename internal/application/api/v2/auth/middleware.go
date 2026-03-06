package auth

import (
	"main/internal/domain/auth"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtClient auth.JwtClient
}

func NewAuthMiddleware(jwtClient auth.JwtClient) *AuthMiddleware {
	return &AuthMiddleware{
		jwtClient: jwtClient,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		token, err := pkg.ParseBearerToken(authHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		payload, err := m.jwtClient.ValidateToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Set("user_id", payload.UserID)
		ctx.Set("is_active", payload.IsActive)

		ctx.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		token, err := pkg.ParseBearerToken(authHeader)
		if err != nil {
			ctx.Next()
			return
		}

		payload, err := m.jwtClient.ValidateToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Set("user_id", payload.UserID)
		ctx.Set("is_active", payload.IsActive)

		ctx.Next()
	}
}

func (m *AuthMiddleware) RequireActive() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isActive := ctx.GetBool("is_active")
		if !isActive {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "ACCOUNT_INACTIVE",
			})
			return
		}

		ctx.Next()
	}
}
