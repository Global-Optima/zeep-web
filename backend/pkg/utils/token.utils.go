package utils

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
)

// HashTokenSHA256 returns a SHA-256 hash of the token as a hex string.
func HashTokenSHA256(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// CompareTokenSHA256 computes the hash of the token and compares it against the stored hash.
func CompareTokenSHA256(storedHash, token string) error {
	computedHash := HashTokenSHA256(token)
	if subtle.ConstantTimeCompare([]byte(storedHash), []byte(computedHash)) == 1 {
		return nil
	}
	return errors.New("token mismatch")
}
