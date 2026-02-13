package service

import "os"

const fallbackUserID = "00000000-0000-0000-0000-000000000001"

func DefaultUserID() string {
	if value := os.Getenv("DEFAULT_USER_ID"); value != "" {
		return value
	}
	return fallbackUserID
}
