package controller

import (
	"testing"

	"github.com/materkov/meme9/web2/store"
	"github.com/stretchr/testify/require"
)

func TestPostRenderer_Render(t *testing.T) {
	r := PostRenderer{
		post: &store.Post{
			ID:     13,
			Text:   "test post",
			UserID: 55,
		},
		user: &store.User{
			ID:   55,
			Name: "User 55 name",
		},
	}

	post := r.Render()
	require.Contains(t, post, "From User 55 name:")
	require.Contains(t, post, "test post")
}
