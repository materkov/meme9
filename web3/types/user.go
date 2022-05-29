package types

import (
	"fmt"
	"github.com/materkov/web3/store"
)

type User struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type UserParams struct {
	Include bool        `json:"include"`
	Name    simpleField `json:"name,omitempty"`
}

func ResolveUser(id int, params UserParams) (*User, error) {
	obj, err := GlobalCachedStore.ObjGet(id)
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %w", err)
	}

	user, ok := obj.(*store.User)
	if !ok {
		return nil, fmt.Errorf("error selecting user: %w", err)
	}

	result := &User{
		Type: "User",
		ID:   fmt.Sprintf("User:%d", user.ID),
	}

	if params.Name.Include {
		result.Name = user.Name
	}

	return result, nil
}
