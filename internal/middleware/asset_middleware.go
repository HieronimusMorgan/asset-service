package middleware

import (
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetMiddleware interface {
	HandlerAsset() gin.HandlerFunc
	HandlerAssetGroup() gin.HandlerFunc
}

type assetMiddleware struct {
	JWTService utils.JWTService
}

func NewAssetMiddleware(jwtService utils.JWTService) AssetMiddleware {
	return assetMiddleware{
		JWTService: jwtService,
	}
}

func (a assetMiddleware) HandlerAsset() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Authorization header is required")
			c.Abort()
			return
		}

		_, err := a.JWTService.ValidateToken(token)
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

		if tokenClaims.Authorized == false {
			response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "You are not authorized to access this resource")
			c.Abort()
			return
		}

		if tokenClaims.Exp < utils.GetCurrentTime() {
			response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "Token has expired")
			c.Abort()
			return
		}

		if !utils.HasAssetResource(tokenClaims.Resource) {
			response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "You are not authorized to access this resource")
			c.Abort()
			return
		}

		c.Set("token", tokenClaims)
		c.Next()
	}
}

func (a assetMiddleware) HandlerAssetGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.SendResponse(c, http.StatusUnauthorized, "Missing token", nil, "Authorization header is required")
			c.Abort()
			return
		}

		_, err := a.JWTService.ValidateToken(token)
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

		if tokenClaims.Authorized == false {
			response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "You are not authorized to access this resource")
			c.Abort()
			return
		}

		if tokenClaims.Exp < utils.GetCurrentTime() {
			response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "Token has expired")
			c.Abort()
			return
		}

		if !utils.HasAssetGroupResource(tokenClaims.Resource) {
			response.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, "You are not authorized to access this resource")
			c.Abort()
			return
		}

		c.Set("token", tokenClaims)
		c.Next()
	}
}
