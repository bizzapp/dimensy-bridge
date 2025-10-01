package utils

import "time"

func NowUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}
