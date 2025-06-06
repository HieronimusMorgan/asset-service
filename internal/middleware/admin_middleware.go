package middleware

import (
	"asset-service/internal/utils/jwt"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AdminMiddleware defines the contract for authentication middleware
type AdminMiddleware interface {
	HandlerAsset() gin.HandlerFunc
}

// adminMiddleware is the struct that implements AdminMiddleware
type adminMiddleware struct {
	JWTService jwt.Service
}

// NewAdminMiddleware initializes authentication middleware
func NewAdminMiddleware(jwtService jwt.Service) AdminMiddleware {
	return adminMiddleware{
		JWTService: jwtService,
	}
}

// Handler returns a middleware function for JWT validation
func (a adminMiddleware) HandlerAsset() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Authorization header is required")
			c.Abort()
			return
		}

		_, err := a.JWTService.ValidateTokenAdmin(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, err.Error())
			c.Abort()
			return
		}

		tokenClaims, err := a.JWTService.ExtractClaims(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token claims", nil, err.Error())
			c.Abort()
			return
		}

		c.Set("token", tokenClaims)

		c.Next()
	}
}
