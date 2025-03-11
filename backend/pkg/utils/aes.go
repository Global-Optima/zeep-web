package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// Struct to store encrypted data
type EncryptedData struct {
	IV      string `json:"iv" binding:"required"`
	Payload string `json:"payload" binding:"required"`
}

// Convert secret key to a 32-byte AES key (SHA-256)
func ConvertSecretKeyToAESKey(secretKey string) []byte {
	hash := sha256.Sum256([]byte(secretKey))
	return hash[:]
}

// Encrypt data using AES-GCM (256-bit)
func EncryptPayload(plainData []byte, secretKey string) (*EncryptedData, error) {
	aesKey := ConvertSecretKeyToAESKey(secretKey)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// **Encrypt and get authentication tag inside ciphertext**
	cipherData := aesGCM.Seal(nil, nonce, plainData, nil)

	return &EncryptedData{
		IV:      base64.StdEncoding.EncodeToString(nonce),
		Payload: base64.StdEncoding.EncodeToString(cipherData),
	}, nil
}

// Decrypt data using AES-GCM (256-bit)
func DecryptPayload(data EncryptedData, secretKey string) ([]byte, error) {
	aesKey := ConvertSecretKeyToAESKey(secretKey)

	iv, err := base64.StdEncoding.DecodeString(data.IV)
	if err != nil {
		return nil, fmt.Errorf("failed to decode IV: %w", err)
	}

	cipherData, err := base64.StdEncoding.DecodeString(data.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// **Decrypt using AES-GCM with authentication tag**
	plainData, err := aesGCM.Open(nil, iv, cipherData, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plainData, nil
}
