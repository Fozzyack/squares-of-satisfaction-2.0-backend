package envutils

import "os"

func GetProduction() bool {
	env := os.Getenv("ENVIRONMENT")
	if env == "prod" || env == "production" {
		return true
	}
	return false
}
