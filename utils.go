package GoTOTP

import (
	"crypto/rand"
	"encoding/base32"
)

// This will generate random bytes and encode them to Base32 (without padding)
func GenerateRandomSecret(length int) (string, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	b32enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	result := b32enc.EncodeToString(secret)
	return result, nil
}
