package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/store"
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

	tokenStr := base64.RawURLEncoding.EncodeToString(tokenBytes)

	h := hmac.New(sha256.New, []byte(store.GlobalConfig.AuthTokenSecret))
	h.Write([]byte(tokenStr))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s.%s", tokenStr, calculatedHash)
}

func ParseAuthToken(tokenStr string) *AuthToken {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 2 {
		return nil
	}

	h := hmac.New(sha256.New, []byte(store.GlobalConfig.AuthTokenSecret))
	h.Write([]byte(parts[0]))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if calculatedHash != parts[1] {
		return nil
	}

	partBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil
	}

	token := AuthToken{}
	err = json.Unmarshal(partBytes, &token)
	if err != nil || token.UserID == 0 {
		return nil
	}

	return &token
}
