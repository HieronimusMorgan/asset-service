package middleware

import (
	"asset-service/internal/models"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Missing token")
			c.Abort()
			return
		}

		claims, err := utils.ExtractClaims(token)
		if err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, err.Error())
			c.Abort()
			return
		}

		var session *models.UserSession
		if err := utils.GetDataFromRedis(utils.UserSession, claims.ClientID, &session); err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid or inactive token", nil, err.Error())
			c.Abort()
			return
		}

		if session.SessionToken != token || !session.IsActive || session.ExpiresAt.Before(time.Now()) {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid or inactive token", nil, "Invalid or inactive token")
			c.Abort()
			return
		}

		if _, err := utils.ValidateToken(token); err != nil {
			response.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
