package globalid

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type ID interface {
	GlobalID()
}

type PostID struct {
	PostID int `json:"postId"`
}

func (PostID) GlobalID() {}

type UserID struct {
	UserID int `json:"userId"`
}

func (UserID) GlobalID() {}

type Stub struct {
}

func (Stub) GlobalID() {}

func Create(id ID) string {
	idBytes, _ := json.Marshal(id)
	prefix := ""

	switch id.(type) {
	case PostID:
		prefix = "PostID"
	case UserID:
		prefix = "UserID"
	case Stub:
		prefix = "Stub"
	}

	idStr := fmt.Sprintf("%s:%s", prefix, idBytes)
	idStr = base64.RawURLEncoding.EncodeToString([]byte(idStr))

	return idStr
}

var ErrIncorrectID = fmt.Errorf("incorrect global id")

func Parse(id string) (ID, error) {
	idBytes, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return nil, ErrIncorrectID
	}

	parts := strings.SplitN(string(idBytes), ":", 2)
	if len(parts) != 2 {
		return nil, nil
	}

	var data ID
	switch parts[0] {
	case "PostID":
		data = &PostID{}
	case "UserID":
		data = &UserID{}
	default:
		return nil, ErrIncorrectID
	}

	err = json.Unmarshal([]byte(parts[1]), data)
	if err != nil {
		return nil, ErrIncorrectID
	}

	return data, nil
}
