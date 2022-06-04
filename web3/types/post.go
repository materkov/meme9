package types

import (
	"fmt"
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
	Inner UserParams `json:"inner"`
}

func ResolveGraphPost(id int, params PostParams) (*Post, error) {
	obj, err := GlobalCachedStore.ObjGet(id)
	if err != nil {
		return nil, fmt.Errorf("error selecting post: %w", err)
	}

	post, ok := obj.(*store.Post)
	if !ok {
		return nil, fmt.Errorf("error selecting post: %w", err)
	}

	result := &Post{
		Type: "Post",
		ID:   fmt.Sprintf("Post:%d", post.ID),
	}

	if params.Text != nil {
		result.Text = post.Text
		if params.Text.MaxLength > 0 && len(result.Text) > params.Text.MaxLength {
			result.Text = result.Text[:params.Text.MaxLength] + "..."
		}
	}

	if params.User != nil {
		result.User, err = ResolveUser(post.UserID, params.User.Inner)
	}

	if params.Date != nil {
		result.Date = post.Date
	}

	return result, nil
}
