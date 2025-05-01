package utils

import (
	response "asset-service/internal/dto/out/assets"
	"errors"
	"time"
)

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

func ToDateOnly(t *time.Time) *response.DateOnly {
	if t == nil {
		return nil
	}
	return (*response.DateOnly)(t)
}
