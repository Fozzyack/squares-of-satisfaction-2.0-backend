package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {

	hash := bcrypt.GenerateFromPassword()
	return hash, nil
}
