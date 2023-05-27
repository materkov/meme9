package api

import (
	"github.com/materkov/meme9/web6/pkg"
	"strconv"
	"time"
)

//go:generate easyjson -all -lower_camel_case -omit_empty api_users.go

type User struct {
	ID   string
	Name string
}

func transformUser(userID int, user *pkg.User) *User {
	result := &User{
		ID: strconv.Itoa(userID),
	}
	if user == nil {
		return result
	}

	result.Name = user.Name

	return result
}

func transformDate(ts int) string {
	return time.Unix(int64(ts), 0).UTC().Format(time.RFC3339)
}
