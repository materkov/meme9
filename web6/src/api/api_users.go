package api

import (
	"errors"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"log"
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
		_, err := store.GlobalStore.GetEdge(viewerID, userID, store.EdgeTypeFollowing)
		if errors.Is(err, store.ErrNoEdge) {
			// Do nothing
		} else if err != nil {
			return nil, err
		} else {
			result.IsFollowing = true
		}
	}

	return result, nil
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

		result[i], err = transformUser(userId, user, v.UserID)
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

	user, err := store.GetUser(v.UserID)
	if err != nil {
		return nil, err
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
		err := store.GlobalStore.DelEdge(v.UserID, targetID, store.EdgeTypeFollowing)
		pkg.LogErr(err)

		err = store.GlobalStore.DelEdge(targetID, v.UserID, store.EdgeTypeFollowedBy)
		pkg.LogErr(err)
	} else {
		_, err := store.GetUser(targetID)
		if errors.Is(err, store.ErrObjectNotFound) {
			return nil, Error("UserNotFound")
		} else if err != nil {
			return nil, err
		}

		err = store.GlobalStore.AddEdge(v.UserID, targetID, store.EdgeTypeFollowing)
		if errors.Is(err, store.ErrDuplicateEdge) {
			// Do nothing
		} else if err != nil {
			return nil, err
		}

		err = store.GlobalStore.AddEdge(targetID, v.UserID, store.EdgeTypeFollowedBy)
		if errors.Is(err, store.ErrDuplicateEdge) {
			// Do nothing
		} else if err != nil {
			return nil, err
		}
	}

	return &Void{}, nil
}
