package auth

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pwByte := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(pwByte, 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func GenerateToken() (string, error) {
	token := make([]byte, 60)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil

}
