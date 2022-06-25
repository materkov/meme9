package types

import (
	"fmt"
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/pkg/globalid"
	"github.com/materkov/web3/store"
)

type Post struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	Text string `json:"text,omitempty"`
	User *User  `json:"user,omitempty"`
	Date int    `json:"date,omitempty"`
}

type PostParams struct {
	Date *simpleField `json:"date"`
	Text *PostText    `json:"text"`
	User *PostUser    `json:"user"`
}

type PostText struct {
	MaxLength int `json:"maxLength,omitempty"`
}

type PostUser struct {
	Inner *UserParams `json:"inner"`
}

func ResolveGraphPost(cachedStore *store.CachedStore, id int, params *PostParams, viewer *pkg.Viewer) (*Post, error) {
	obj, err := cachedStore.ObjGet(id)
	if err != nil {
		return nil, fmt.Errorf("error selecting post: %w", err)
	}

	post, ok := obj.(*store.Post)
	if !ok {
		return nil, fmt.Errorf("post not found")
	}

	result := &Post{
		Type: "Post",
		ID:   globalid.Create(globalid.PostID{PostID: post.ID}),
	}

	if params == nil {
		return result, err
	}

	if params.Text != nil {
		result.Text = post.Text
		if params.Text.MaxLength > 0 && len(result.Text) > params.Text.MaxLength {
			result.Text = result.Text[:params.Text.MaxLength] + "..."
		}
	}

	if params.User != nil {
		result.User, err = ResolveUser(cachedStore, post.UserID, params.User.Inner, viewer)
	}

	if params.Date != nil {
		result.Date = post.Date
	}

	return result, nil
}
