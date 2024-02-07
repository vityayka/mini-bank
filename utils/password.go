package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate a hash; %v", err)
	}
	return string(hash), nil
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
