package utils

import (
	"fmt"
	"time"
)

func StringTimeToISO(input string) (time.Time, error) {
	layout := "15:04"
	parsedTime, err := time.Parse(layout, input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format, expected HH:MM: %w", err)
	}

	formattedTime := time.Date(1970, 1, 1, parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.UTC)
	return formattedTime, nil
}
