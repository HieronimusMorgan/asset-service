package utils

import (
	response "asset-service/internal/dto/out/assets"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ValidationTrimSpace(s string) string {
	trim := strings.TrimSpace(s)
	trim = strings.Join(strings.Fields(trim), " ") // Remove extra spaces
	return trim
}

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsername.MatchString(username) {
		return errors.New("username can only contain alphanumeric characters and underscores")
	}

	return nil
}

func ConvertToUint(input string) (uint, error) {
	parsed, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid uint price: %w", err)
	}
	return uint(parsed), nil
}

func ParseDate(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, errors.New("error parsing date: " + err.Error())
	}
	return parsedDate, nil
}

func GetToday() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	location := now.Location()
	return time.Date(year, month, day, 0, 0, 0, 0, location)
}

func GenerateInviteToken() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 32-character hex string
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ParseOptionalDate(str *string) (*time.Time, error) {
	if str == nil {
		return nil, nil
	}
	parsedDate, err := time.Parse("2006-01-02", *str)
	if err != nil {
		return nil, err
	}
	return &parsedDate, nil
}

func CalculateNextDueDate(date *time.Time, days *int) (*time.Time, error) {
	if date == nil || days == nil {
		return nil, nil
	}
	nextDueDate := date.AddDate(0, 0, *days)
	return &nextDueDate, nil
}
func NullableStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ToDateOnly(t *time.Time) *response.DateOnly {
	if t == nil {
		return nil
	}
	return (*response.DateOnly)(t)
}

func GetOptionalString(context *gin.Context, field string) *string {
	val := context.PostForm(field)
	if val == "" {
		return nil
	}
	return &val
}

func ParseFormInt(context *gin.Context, field string) int {
	val := context.PostForm(field)
	if val == "" {
		return 0
	}
	intVal, _ := strconv.Atoi(val)
	return intVal
}

func ParseFormUint(context *gin.Context, field string) uint {
	val := context.PostForm(field)
	if val == "" {
		return 0
	}
	uintVal, _ := strconv.ParseUint(val, 10, 32)
	return uint(uintVal)
}

func ParseFormFloat(context *gin.Context, field string) float64 {
	val := context.PostForm(field)
	if val == "" {
		return 0.0
	}
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}

// CheckCredentialKey checks if the provided credential key matches the one stored in Redis for the given client ID.
func CheckCredentialKey(redis RedisService, credential, clientID string) error {
	credentialKeyMap := struct {
		CredentialKey string `json:"credential_key"`
	}{}
	err := redis.GetData(CredentialKey, clientID, &credentialKeyMap)

	log.Info().Str("checkCredential", credentialKeyMap.CredentialKey).Msg("Credential key retrieved from Redis")
	log.Info().Str("credential", credential).Msg("Credential key retrieved from Redis")
	if err != nil {
		log.Error().Str("credentialKey", credential).Err(err).Msg("Failed to retrieve credential key from Redis")
		return err
	}

	if credentialKeyMap.CredentialKey != credential {
		return errors.New("credential key not matched")
	}
	_ = redis.DeleteData(CredentialKey, clientID)

	return nil
}
