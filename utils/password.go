package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
)

// GenerateRandomSalt creates len bytes of salt
func GenerateRandomSalt(len int) ([]byte, error) {
	salt := make([]byte, len)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

// HashPassword combines password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed password
// as a hex string
func HashPassword(plainTextPassword string, salt []byte) string {
	passwordBytes := []byte(plainTextPassword)

	// Create sha-512 hasher
	sha512Hasher := sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	hashedPasswordBytes := sha512Hasher.Sum(nil)

	// Convert the hashed password to a hex string
	hashedPasswordHex := hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

func IsPasswordsMatch(hashedPassword, plainTextPassword string, salt []byte) bool {
	return hashedPassword == HashPassword(plainTextPassword, salt)
}
