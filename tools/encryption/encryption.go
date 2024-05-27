package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

// IsValid - validates if two passwords match.
func IsValid(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err
}

// Encrypt - Generate hashed password out of clear text password.
func Encrypt(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
