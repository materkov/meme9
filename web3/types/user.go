package types

import (
	"fmt"
	"github.com/materkov/web3/pkg/globalid"
	"github.com/materkov/web3/store"
)

type User struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`

	Name  string               `json:"name,omitempty"`
	Posts *UserPostsConnection `json:"posts,omitempty"`
}

type UserParams struct {
	Name  *simpleField               `json:"name,omitempty"`
	Posts *UserPostsConnectionFields `json:"posts,omitempty"`
}

func ResolveUser(cachedStore *store.CachedStore, id int, params *UserParams) (*User, error) {
	obj, err := cachedStore.ObjGet(id)
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %w", err)
	}

	user, ok := obj.(*store.User)
	if !ok {
		return nil, fmt.Errorf("error selecting user: %w", err)
	}

	result := &User{
		Type: "User",
		ID:   globalid.Create(globalid.UserID{UserID: user.ID}),
	}

	if params.Name != nil {
		result.Name = user.Name
	}

	if params.Posts != nil {
		result.Posts, _ = ResolveUserPostsConnection(cachedStore, user.ID, params.Posts)
	}

	return result, nil
}
