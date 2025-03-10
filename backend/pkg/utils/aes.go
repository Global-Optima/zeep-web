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

type EncryptedData struct {
	IV      string `json:"iv" binding:"required"`
	Payload string `json:"payload" binding:"required"`
}

func ConvertSecretKeyToAESKey(secretKey string) []byte {
	hash := sha256.Sum256([]byte(secretKey)) // Hashing ensures a fixed-length 32-byte key
	return hash[:]
}

func EncryptPayload(plainData []byte, secretKey string) (*EncryptedData, error) {
	aesKey := ConvertSecretKeyToAESKey(secretKey) // Convert secretKey to 32-byte AES key

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

	cipherData := aesGCM.Seal(nil, nonce, plainData, nil)

	return &EncryptedData{
		IV:      base64.StdEncoding.EncodeToString(nonce),
		Payload: base64.StdEncoding.EncodeToString(cipherData),
	}, nil
}

func DecryptPayload(data EncryptedData, secretKey string) ([]byte, error) {
	aesKey := ConvertSecretKeyToAESKey(secretKey) // Convert secretKey to 32-byte AES key

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

	plainData, err := aesGCM.Open(nil, iv, cipherData, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plainData, nil
}
