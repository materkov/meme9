package types

import (
	"fmt"
	"github.com/materkov/web3/pkg/globalid"
	"github.com/materkov/web3/store"
)

type User struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`

	Avatar string               `json:"avatar,omitempty"`
	Name   string               `json:"name,omitempty"`
	Posts  *UserPostsConnection `json:"posts,omitempty"`
}

type UserParams struct {
	Name   *simpleField               `json:"name,omitempty"`
	Posts  *UserPostsConnectionFields `json:"posts,omitempty"`
	Avatar *simpleField               `json:"avatar,omitempty"`
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

	if params.Avatar != nil {
		fileHash := user.AvatarFile
		if fileHash == "" {
			fileHash = "dbb7f7e5b2658593b648328c3bdc95ad0253e65e816d061d789de09f81663a5d"
		}

		result.Avatar = fmt.Sprintf("https://689809.selcdn.ru/meme-files/avatars/%s", fileHash)
	}

	return result, nil
}
