package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(rawToken string) string {
	sum := sha256.Sum224([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
