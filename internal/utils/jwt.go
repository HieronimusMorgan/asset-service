package utils

import (
	"asset-service/config"
	"asset-service/package/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var internalSecretKey = []byte(os.Getenv("JWT_SECRET"))

func ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	secret := config.GetJWTSecret() // ✅ Load secret dynamically

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil // ✅ Use dynamically loaded secret
	})

	if err != nil {
		return nil, err
	}

	// Extract claims and validate
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
	}

	return &claims, nil
}

func ValidateTokenAdmin(tokenString string) (*jwt.MapClaims, error) {
	secret := config.GetJWTSecret() // ✅ Load secret dynamically

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil // ✅ Use dynamically loaded secret
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
	}

	if role, ok := claims["role"].(string); ok {
		if strings.EqualFold(role, "Admin") || strings.EqualFold(role, "Super Admin") {
			return &claims, nil
		}
		return nil, errors.New("user is not an Admin")
	}

	return nil, errors.New("role not found in token claims")
}

func ExtractClaims(tokenString string) (*TokenClaims, error) {
	secret := config.GetJWTSecret() // ✅ Load secret dynamically

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil // ✅ Use dynamically loaded secret
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	tc := &TokenClaims{}

	if authorized, ok := claims["authorized"].(bool); ok {
		tc.Authorized = authorized
	}

	if accessUUID, ok := claims["access_uuid"].(string); ok {
		tc.AccessUUID = accessUUID
	}

	if exp, ok := claims["exp"].(float64); ok {
		tc.Exp = int64(exp)
	}

	if userID, ok := claims["user_id"].(float64); ok {
		tc.UserID = uint(userID)
	}

	if clientID, ok := claims["client_id"].(string); ok {
		tc.ClientID = clientID
	}

	if role, ok := claims["role_id"].(float64); ok {
		tc.RoleID = uint(role)
	}

	return tc, nil
}

type TokenClaims struct {
	Authorized bool   `json:"authorized"`
	AccessUUID string `json:"access_uuid"`
	UserID     uint   `json:"user_id"`
	ClientID   string `json:"client_id"`
	RoleID     uint   `json:"role_id"`
	Exp        int64  `json:"exp"`
}

type InternalClaims struct {
	Service string `json:"service"` // Service name
	jwt.RegisteredClaims
}

func ExtractClaimsResponse(context *gin.Context) (TokenClaims, error) {
	token, err := ExtractClaims(context.GetHeader("Authorization"))
	if err != nil {
		response.SendResponse(context, 401, "Unauthorized", nil, err.Error())
	}
	return *token, err
}
