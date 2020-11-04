package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type AddPost struct {
	Store *store.Store
}

func (a *AddPost) Handle(viewer *api.Viewer, req *pb.AddPostRequest) (*pb.AddPostRenderer, error) {
	postID, err := a.Store.GenerateNodeID()
	if err != nil {
		return nil, fmt.Errorf("error generating node id: %w", err)
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
		return nil, fmt.Errorf("error adding post: %w", err)
	}

	err = a.Store.AddToFeed(post.ID)
	if err != nil {
		log.Printf("[ERROR] Error adding post to feed: %s", err)
	}

	renderer := &pb.AddPostRenderer{
		Id:   strconv.Itoa(post.ID),
		Text: post.Text,
	}

	return renderer, nil
}
