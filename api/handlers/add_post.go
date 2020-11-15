package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
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

	if viewer.User == nil || !viewer.CSRFValidated {
		return nil, api.NewError("NOT_AUTHORIZED", "User not authorized")
	}

	text := strings.TrimSpace(req.Text)
	if text == "" {
		return nil, api.NewError("TEXT_EMPTY", "Text is empty")
	} else if len(text) > 1000 {
		return nil, api.NewError("TEXT_TOO_LONG", "Text is too long")
	}

	post := store.Post{
		ID:        postID,
		Text:      text,
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
		Id:          strconv.Itoa(post.ID),
		Text:        post.Text,
		SuccessText: "Пост добавлен",
		PostUrl:     fmt.Sprintf("/posts/%d", post.ID),
	}

	return renderer, nil
}
