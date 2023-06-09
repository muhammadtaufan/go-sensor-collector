package pkg

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func ParseDateWithFallback(dateStr, fallbackFormat, defaultTime, timezone string) (*time.Time, error) {
	location, _ := time.LoadLocation(timezone)

	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		parsedDate, err := time.ParseInLocation(format, dateStr, location)
		if err == nil {
			if format == fallbackFormat {
				parsedDate, _ = time.ParseInLocation("2006-01-02 15:04:05",
					parsedDate.Format("2006-01-02")+defaultTime, location)
			}
			return &parsedDate, nil
		}
	}

	return nil, fmt.Errorf("invalid date format")
}

func ValidateHashPassword(dbPassword, requestPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(requestPassword))
	if err != nil {
		return false, fmt.Errorf("invalid password")
	}

	return true, nil
}
