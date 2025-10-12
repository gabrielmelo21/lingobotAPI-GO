package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateReferralCode() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
