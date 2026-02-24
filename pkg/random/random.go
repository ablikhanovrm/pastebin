package random

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRefreshToken(bytesLen int) (string, error) {
	b := make([]byte, bytesLen)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
