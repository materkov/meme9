package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type AuthToken struct {
	UserID int `json:"userId"`
}

func (a *AuthToken) ToString() string {
	tokenBytes, err := json.Marshal(a)
	if err != nil {
		return ""
	}

	h := hmac.New(sha256.New, []byte(GlobalConfig.AuthTokenSecret))
	h.Write(tokenBytes)
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s.%s", tokenBytes, calculatedHash)
}

func ParseAuthToken(tokenStr string) *AuthToken {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 2 {
		return nil
	}

	h := hmac.New(sha256.New, []byte(GlobalConfig.AuthTokenSecret))
	h.Write([]byte(parts[0]))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if calculatedHash != parts[1] {
		return nil
	}

	token := AuthToken{}
	err := json.Unmarshal([]byte(parts[0]), &token)
	if err != nil || token.UserID == 0 {
		return nil
	}

	return &token
}
