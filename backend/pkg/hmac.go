package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

func VerifyHMAC(payload, hmacSignature string, sharedKey []byte) error {
	signatureBytes, err := hex.DecodeString(hmacSignature)
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}

	// Compute HMAC of payload using shared key
	mac := hmac.New(sha256.New, sharedKey)
	mac.Write([]byte(payload))
	expectedMAC := mac.Sum(nil)

	// Securely compare the provided HMAC with the expected one
	if !hmac.Equal(signatureBytes, expectedMAC) {
		return errors.New("HMAC signature mismatch")
	}

	return nil // HMAC is valid
}

func DecryptPayload(encryptedBase64 string, sharedKey []byte) ([]byte, error) {
	// Decode base64-encoded encrypted data
	cipherData, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, fmt.Errorf("decode base64 failed: %w", err)
	}

	// Ensure key length is 32 bytes for AES-256
	if len(sharedKey) != 32 {
		return nil, fmt.Errorf("invalid key length: must be 32 bytes")
	}

	// Extract nonce and ciphertext
	if len(cipherData) < 12 {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := cipherData[:12], cipherData[12:]

	// Create AES cipher block
	block, err := aes.NewCipher(sharedKey)
	if err != nil {
		return nil, err
	}

	// Create AES-GCM decryption instance
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt data
	plainData, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plainData, nil
}

// EncryptPayload encrypts data using AES-GCM and returns a base64-encoded string.
func EncryptPayload(plainData []byte, sharedKey []byte) (string, error) {
	// Ensure key length is 32 bytes for AES-256
	if len(sharedKey) != 32 {
		return "", fmt.Errorf("invalid key length: must be 32 bytes")
	}

	// Create AES cipher block
	block, err := aes.NewCipher(sharedKey)
	if err != nil {
		return "", err
	}

	// Create AES-GCM encryption instance
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a random nonce (12 bytes)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt data
	cipherData := aesGCM.Seal(nil, nonce, plainData, nil)

	// Prepend nonce to ciphertext and return base64-encoded result
	encryptedBase64 := base64.StdEncoding.EncodeToString(append(nonce, cipherData...))
	return encryptedBase64, nil
}

// GenerateHMAC generates an HMAC signature for a given payload.
func GenerateHMAC(payload string, sharedKey []byte) string {
	mac := hmac.New(sha256.New, sharedKey)
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
