package api

import (
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"strconv"
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
	IsFollowing bool   `json:"isFollowing,omitempty"`
}

func transformUser(userID int, user *store.User, viewerID int) (*User, error) {
	result := &User{
		ID: strconv.Itoa(userID),
	}
	if user == nil {
		return result, nil
	}

	result.Name = user.Name
	result.Status = user.Status

	if viewerID != 0 {
		isFollowing, err := store2.GlobalStore.Subs.CheckFollowing(viewerID, []int{userID})
		if err != nil {
			return nil, err
		} else {
			result.IsFollowing = isFollowing[userID]
		}
	}

	return result, nil
}

type UsersListReq struct {
	UserIds []string `json:"userIds"`
}

func (*API) usersList(v *Viewer, r *UsersListReq) ([]*User, error) {
	userIds := utils.IdsToInts(r.UserIds)
	users, err := store2.GlobalStore.Users.Get(userIds)
	if err != nil {
		return nil, err
	}

	result := make([]*User, len(r.UserIds))
	for i, userIdStr := range r.UserIds {
		userId, _ := strconv.Atoi(userIdStr)

		result[i], err = transformUser(userId, users[userId], v.UserID)
		pkg.LogErr(err)
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

	users, err := store2.GlobalStore.Users.Get([]int{v.UserID})
	if err != nil {
		return nil, err
	}

	user := users[v.UserID]
	if user == nil {
		return nil, fmt.Errorf("cannot find viewer")
	}

	user.Status = r.Status

	err = store2.GlobalStore.Nodes.Update(user.ID, user)
	if err != nil {
		return nil, err
	}

	return &Void{}, nil
}

type SubscribeAction string

const (
	Follow   SubscribeAction = "FOLLOW"
	Unfollow SubscribeAction = "UNFOLLOW"
)

type UsersFollow struct {
	TargetID string          `json:"targetId"`
	Action   SubscribeAction `json:"action"`
}

func (*API) usersFollow(v *Viewer, r *UsersFollow) (*Void, error) {
	if v.UserID == 0 {
		return nil, Error("NotAuthorized")
	}

	targetID, _ := strconv.Atoi(r.TargetID)
	if targetID <= 0 {
		return nil, Error("InvalidTarget")
	}

	if r.Action == Unfollow {
		err := store2.GlobalStore.Subs.Follow(v.UserID, targetID)
		pkg.LogErr(err)
	} else {
		users, err := store2.GlobalStore.Users.Get([]int{targetID})
		if err != nil {
			return nil, err
		} else if users[targetID] == nil {
			return nil, Error("UserNotFound")
		}

		err = store2.GlobalStore.Subs.Follow(v.UserID, targetID)
		if err != nil {
			return nil, err
		}
	}

	return &Void{}, nil
}
