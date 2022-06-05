package types

import (
	"fmt"
	"github.com/materkov/web3/store"
)

type User struct {
	Type  string  `json:"type,omitempty"`
	ID    string  `json:"id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Posts []*Post `json:"posts,omitempty"`
}

type UserParams struct {
	Name  *simpleField `json:"name,omitempty"`
	Posts *UserPosts   `json:"posts,omitempty"`
}

type UserPosts struct {
	Inner PostParams `json:"inner"`
}

func ResolveUser(cachedStore *store.CachedStore, id int, params UserParams) (*User, error) {
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
		ID:   fmt.Sprintf("User:%d", user.ID),
	}

	if params.Name != nil {
		result.Name = user.Name
	}

	if params.Posts != nil {
		postIds, _ := cachedStore.Store.ListGet(user.ID, store.ListPosted)
		for _, postID := range postIds {
			cachedStore.Need(postID)
		}

		for _, postID := range postIds {
			post, _ := ResolveGraphPost(cachedStore, postID, params.Posts.Inner)
			result.Posts = append(result.Posts, post)
		}
	}

	return result, nil
}
