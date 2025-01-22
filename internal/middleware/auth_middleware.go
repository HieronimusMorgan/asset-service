package middleware

import (
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Missing token")
			c.Abort()
			return
		}

		_, err := utils.ValidateTokenAdmin(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Admin Access", nil, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
