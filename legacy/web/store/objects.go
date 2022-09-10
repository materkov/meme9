package store

import (
	"strconv"
	"strings"
)

type StoredObject struct {
	ID int `json:"ID"`

	APILog  *APILog  `json:"APILog,omitempty"`
	Photo   *Photo   `json:"Photo,omitempty"`
	Comment *Comment `json:"Comment,omitempty"`
	User    *User    `json:"User,omitempty"`
	Token   *Token   `json:"Token,omitempty"`
	Post    *Post    `json:"Post,omitempty"`
}

type APILog struct {
	ID       int
	UserID   int
	Method   string
	Request  string
	Response string
}

type Photo struct {
	ID     int
	UserID int
	Path   string
}

type User struct {
	ID       int
	VkID     int
	Name     string
	VkAvatar string
}

type Comment struct {
	ID     int
	PostID int
	UserID int
	Text   string
	Date   int
}

type Token struct {
	ID     int
	Token  string
	UserID int
}

type Post struct {
	ID      int
	UserID  int
	Date    int
	Text    string
	PhotoID int
}

func GetIdFromToken(token string) int {
	parts := strings.Split(token, "-")
	if len(parts) < 2 {
		return 0
	}
	tokenID, _ := strconv.Atoi(parts[0])
	if tokenID <= 0 {
		return 0
	}

	return tokenID
}
