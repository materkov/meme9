package store

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

//go:generate msgp -tests=false

type Token struct {
	ID       int
	IssuedAt int
	Token    string
	Type     int
	UserID   int
}

func (t *Token) GenerateRandomToken() {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 64)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	t.Token = fmt.Sprintf("%d-%d-%s", t.ID, t.UserID, b)
}

func GetNodeIDFromToken(token string) (int, error) {
	parts := strings.Split(token, "-")
	if len(parts) < 2 {
		return 0, fmt.Errorf("incorrect token")
	}

	nodeID, _ := strconv.Atoi(parts[0])
	if nodeID <= 0 {
		return 0, fmt.Errorf("incorrect token")
	}

	return nodeID, nil
}

const (
	TokenTypeCookie  = 1
	TokenTypeRegular = 2
)

type User struct {
	ID           int
	Name         string
	PasswordHash string
	VkID         int
}

type Post struct {
	ID        int
	Text      string
	UserID    int
	Date      int
	UserAgent string
}
