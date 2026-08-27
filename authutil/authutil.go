package authutil

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func KeyGenerator(prefix string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("Error reading random bytes: %w", err)
	}
	key := hex.EncodeToString(b)
	return fmt.Sprintf("%s-%s", prefix, key), nil
}

func HashApiKeys(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
