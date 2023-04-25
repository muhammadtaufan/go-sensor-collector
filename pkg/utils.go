package pkg

import (
	"fmt"
	"time"
)

func ParseDateWithFallback(dateStr, fallbackFormat, defaultTime string) (*time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		parsedDate, err := time.Parse(format, dateStr)
		if err == nil {
			if format == fallbackFormat {
				parsedDate, _ = time.Parse("2006-01-02 15:04:05", parsedDate.Format("2006-01-02")+defaultTime)
			}
			return &parsedDate, nil
		}
	}

	return nil, fmt.Errorf("invalid date format")
}
