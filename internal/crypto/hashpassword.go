package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost represents number of loops within bcrypt
	BcryptCost = 7
)

// HashPassword hashes a password using bcrypt algorithm
func HashPassword(password string) (string, error) {
	// bcrpt password and add to account
	passBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)

	if err != nil {
		return "", err
	}

	passString := string(passBytes)
	return passString, nil
}

// ValidatePassword takes a password and hash string and returns true
// if they validate or false if they do not. An error is returned on any
// other condition
func ValidatePassword(password, hash string) bool {
	bPass, bHash := []byte(password), []byte(hash)

	err := bcrypt.CompareHashAndPassword(bHash, bPass)
	if err != nil {
		return false
	}

	return true
}
