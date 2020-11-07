package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

// TODO: move to redis
var hashKey = []byte("]G4<w}t>;EZA*erX")

func GenerateCSRFToken(viewerID int) string {
	mac := hmac.New(sha256.New, hashKey)
	return fmt.Sprintf("%d-%x", viewerID, mac.Sum([]byte(strconv.Itoa(viewerID))))
}

func ValidateCSRFToken(viewerID int, token string) bool {
	parts := strings.Split(token, "-")
	if len(parts) != 2 {
		return false
	}

	tokenUserID, _ := strconv.Atoi(parts[0])
	if viewerID != tokenUserID {
		return false
	}

	validToken := GenerateCSRFToken(viewerID)
	return validToken == token
}
