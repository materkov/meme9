package pkg

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const (
	GlobalIDUser                = "User"
	GlobalIDPost                = "Post"
	GlobalIDUserPostsConnection = "UserPostsConnection"
)

var allObjectTypes = map[string]string{
	GlobalIDPost:                GlobalIDPost,
	GlobalIDUser:                GlobalIDUser,
	GlobalIDUserPostsConnection: GlobalIDUserPostsConnection,
}

func GetGlobalID(objectType string, objectID int) string {
	globalID := fmt.Sprintf("%s:%d", objectType, objectID)

	return base64.RawURLEncoding.EncodeToString([]byte(globalID))
}

var ErrInvalidGlobalID = fmt.Errorf("invalid global id")

func ParseGlobalID(objectID string) (string, int, error) {
	globalIdBytes, err := base64.RawURLEncoding.DecodeString(objectID)
	if err != nil {
		return "", 0, ErrInvalidGlobalID
	}

	parts := strings.Split(string(globalIdBytes), ":")
	if len(parts) != 2 {
		return "", 0, ErrInvalidGlobalID
	}

	objectType := allObjectTypes[parts[0]]
	if objectType == "" {
		return "", 0, ErrInvalidGlobalID
	}

	objectIDInt, _ := strconv.Atoi(parts[1])
	if objectIDInt <= 0 {
		return "", 0, ErrInvalidGlobalID
	}

	return objectType, objectIDInt, nil
}
