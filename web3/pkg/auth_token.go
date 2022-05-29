package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type AuthToken struct {
	IssuedAt int
	UserID   int
}

func (a *AuthToken) ToString() string {
	part1 := []byte(base64.RawURLEncoding.EncodeToString([]byte(
		`{"alg": "HS256", "typ": "JWT"}`,
	)))

	part2, _ := json.Marshal(a)
	part2 = []byte(base64.RawURLEncoding.EncodeToString(part2))

	h := hmac.New(sha256.New, []byte(GlobalConfig.AuthTokenSecret))
	h.Write(part1)
	h.Write([]byte("."))
	h.Write(part2)

	part3 := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s.%s.%s", part1, part2, part3)
}

var ErrIncorrectToken = fmt.Errorf("incorrect token")

func ParseAuthToken(tokenStr string) (AuthToken, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return AuthToken{}, ErrIncorrectToken
	}

	h := hmac.New(sha256.New, []byte(GlobalConfig.AuthTokenSecret))
	h.Write([]byte(parts[0] + "." + parts[1]))

	if parts[2] != hex.EncodeToString(h.Sum(nil)) {
		return AuthToken{}, ErrIncorrectToken
	}

	part1Bytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return AuthToken{}, ErrIncorrectToken
	}

	token := AuthToken{}
	err = json.Unmarshal(part1Bytes, &token)
	if err != nil {
		return AuthToken{}, ErrIncorrectToken
	}

	if int(time.Now().Unix())-token.IssuedAt > 3600 {
		return AuthToken{}, ErrIncorrectToken
	}

	return token, nil
}
