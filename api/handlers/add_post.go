package handlers

import (
	"strconv"
	"time"

	"github.com/materkov/meme9/api/api"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type AddPost struct {
	Store *store.Store
}

func (a *AddPost) Handle(viewer *api.Viewer, req *login.AddPostRequest) *login.AnyRenderer {
	postID, err := a.Store.GenerateNodeID()
	if err != nil {
		return &login.AnyRenderer{Renderer: &login.AnyRenderer_ErrorRenderer{
			ErrorRenderer: &login.ErrorRenderer{
				DisplayText: "eer",
			},
		}}
	}

	post := store.Post{
		ID:        postID,
		Text:      req.Text,
		UserID:    viewer.User.ID,
		Date:      int(time.Now().Unix()),
		UserAgent: viewer.UserAgent,
	}

	err = a.Store.AddPost(&post)
	if err != nil {
		return &login.AnyRenderer{Renderer: &login.AnyRenderer_ErrorRenderer{
			ErrorRenderer: &login.ErrorRenderer{
				DisplayText: "eer",
			},
		}}
	}

	_ = a.Store.AddToFeed(post.ID)

	return &login.AnyRenderer{Renderer: &login.AnyRenderer_AddPostRenderer{
		AddPostRenderer: &login.AddPostRenderer{
			Id:   strconv.Itoa(post.ID),
			Text: post.Text,
		},
	}}
}
