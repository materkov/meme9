package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
)

func GenerateToken(tokenSecret string, viewerID int) string {
	mac := hmac.New(sha256.New, []byte(tokenSecret))
	_, _ = mac.Write([]byte(strconv.Itoa(viewerID)))
	return fmt.Sprintf("%d-%x", viewerID, mac.Sum(nil))
}

func ValidateToken(tokenSecret string, viewerID int, token string) bool {
	parts := strings.Split(token, "-")
	if len(parts) != 2 {
		return false
	}

	tokenUserID, _ := strconv.Atoi(parts[0])
	if viewerID != tokenUserID {
		return false
	}

	validToken := GenerateToken(tokenSecret, viewerID)
	return validToken == token
}
