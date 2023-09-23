package api

import (
	"github.com/materkov/meme9/web6/src/store"
	"log"
	"strconv"
)

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

func transformUser(userID int, user *store.User) *User {
	result := &User{
		ID: strconv.Itoa(userID),
	}
	if user == nil {
		return result
	}

	result.Name = user.Name
	result.Status = user.Status

	return result
}

type UsersListReq struct {
	UserIds []string `json:"userIds"`
}

func (*API) usersList(v *Viewer, r *UsersListReq) ([]*User, error) {
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

type UsersSetStatus struct {
	Status string `json:"status"`
}

func (*API) usersSetStatus(v *Viewer, r *UsersSetStatus) (*Void, error) {
	if v.UserID == 0 {
		return nil, Error("NotAuthorized")
	}
	if len(r.Status) > 100 {
		return nil, Error("StatusTooLong")
	}

	user, err := store.GetUser(v.UserID)
	if err != nil {
		return nil, err
	}

	user.Status = r.Status

	err = store.UpdateObject(user, user.ID)
	if err != nil {
		return nil, err
	}

	return &Void{}, nil
}
