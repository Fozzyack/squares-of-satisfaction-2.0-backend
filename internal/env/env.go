package env

import (
	"fmt"
	"os"
)

func GetProduction() bool {
	env := os.Getenv("ENVIRONMENT")
	if env == "prod" || env == "production" || env == "PROD" || env == "PRODUCTION" {
		return true
	}
	return false
}

func GetDbConnectionString() (string, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return "", fmt.Errorf("Could not Read DATABASE_URL is it set as an ENVIRONMENT variable?")
	}
	return connStr, nil
}
