package api

import (
	"github.com/materkov/meme9/web6/src/store"
	"log"
	"strconv"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func transformUser(userID int, user *store.User) *User {
	result := &User{
		ID: strconv.Itoa(userID),
	}
	if user == nil {
		return result
	}

	result.Name = user.Name

	return result
}

type UsersListReq struct {
	UserIds []string `json:"userIds"`
}

func (*API) usersList(v *Viewer, r *UsersListReq) (interface{}, error) {
	result := make([]*User, len(r.UserIds))
	for i, userIdStr := range r.UserIds {
		userId, _ := strconv.Atoi(userIdStr)
		user, err := store.GetUser(userId)
		if err != nil {
			log.Printf("[ERROR] Error loading user: %s", err)
		}

		result[i] = transformUser(userId, user)
	}

	return result, nil
}
