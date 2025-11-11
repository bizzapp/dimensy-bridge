package utils

import "time"

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseBirthDate parses date string in format "1994-08-25" to *time.Time
func ParseBirthDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}
