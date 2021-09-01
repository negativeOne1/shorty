package random

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

func GetToken(length int) (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return strings.ToLower(base32.StdEncoding.EncodeToString(randomBytes)[:length]), nil
}
