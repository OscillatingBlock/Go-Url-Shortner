package app

import (
	"math/rand"
	"strings"
)

var (
	SHORT_CODE_LEN int    = 7
	ALLOWED_CHARS  string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateRandomString(length int) string {
	var randomStr strings.Builder
	randomStr.Grow(length)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(ALLOWED_CHARS))
		randomStr.WriteByte(ALLOWED_CHARS[randomIndex])
	}
	return randomStr.String()
}
