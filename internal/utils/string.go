package utils

import (
	"errors"
	"fmt"
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

func NilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
