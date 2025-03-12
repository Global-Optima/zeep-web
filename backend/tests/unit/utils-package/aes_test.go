package utils

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type TestMessage struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func TestEncryptDecryptJSON(t *testing.T) {
	// Given test key (Base64-encoded)
	secretKey := "AJmQlXMjVrcKXxJnkeS5Q2BaJ/9EtIKCqhhgl+RopQU="

	// Decode the Base64 key to get a valid string passphrase
	keyBytes, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		t.Fatalf("Failed to decode Base64 secret key: %v", err)
	}

	// Convert to string (as used in ConvertSecretKeyToAESKey)
	secretKeyString := string(keyBytes)

	// Sample JSON struct
	originalData := TestMessage{
		ID:      123,
		Message: "Hello, this is a JSON test!",
	}

	// Marshal JSON to bytes
	jsonData, err := json.Marshal(originalData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Encrypt the JSON
	encrypted, err := utils.EncryptPayload(jsonData, secretKeyString)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Ensure IV and Payload are non-empty
	if encrypted.IV == "" || encrypted.Payload == "" {
		t.Fatalf("Encryption produced empty IV or Payload")
	}

	// Decrypt the JSON
	decrypted, err := utils.DecryptPayload(*encrypted, secretKeyString)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Unmarshal decrypted data back into struct
	var decryptedData TestMessage
	err = json.Unmarshal(decrypted, &decryptedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal decrypted JSON: %v", err)
	}

	// Assert decrypted struct matches original
	if decryptedData.ID != originalData.ID || decryptedData.Message != originalData.Message {
		t.Fatalf("Decryption output mismatch: got %+v, expected %+v", decryptedData, originalData)
	}
}
