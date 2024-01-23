package store

import (
	"fmt"
	"math/rand"
)

type Token struct {
	ID     int
	UserID int
	Date   int
	Token  string
}

func GenerateToken(userID int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const tokenLength = 32

	randPart := make([]byte, tokenLength)
	for i := range randPart {
		randPart[i] = charset[rand.Intn(len(charset))]
	}

	return fmt.Sprintf("t-%d-%s", userID, randPart)
}
