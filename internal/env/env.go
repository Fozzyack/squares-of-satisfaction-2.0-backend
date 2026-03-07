package env

import "os"

func GetProduction() bool {
	env := os.Getenv("ENVIRONMENT")
	if env == "prod" || env == "production" || env == "PROD" || env == "PRODUCTION" {
		return true
	}
	return false
}
