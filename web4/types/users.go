package types

import (
	"fmt"
	"github.com/materkov/meme9/web4/store"
	"strconv"
)

type User struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Href string `json:"href,omitempty"`
}

func usersList(ids []int) []*User {
	usersMap := map[int]store.User{}
	for _, postID := range ids {
		obj := store.User{}
		err := getObject(postID, &obj)
		if err == nil {
			usersMap[obj.ID] = obj
		}
	}

	results := make([]*User, len(ids))
	for i, userID := range ids {
		result := &User{
			ID:   strconv.Itoa(userID),
			Href: fmt.Sprintf("/users/%d", userID),
		}
		results[i] = result

		user, ok := usersMap[userID]
		if !ok {
			continue
		}

		result.Name = user.Name
	}

	return results
}
