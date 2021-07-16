package app

import (
	"context"
	"log"

	"github.com/materkov/meme9/web/pb"
)

type Feed struct {
	App *App
}

func (f *Feed) GetHeader(ctx context.Context, _ *pb.FeedGetHeaderRequest) (*pb.FeedGetHeaderResponse, error) {
	viewer := GetViewerFromContext(ctx)

	headerRenderer := pb.HeaderRenderer{
		MainUrl:   "/",
		LoginUrl:  "/login",
		LogoutUrl: "/logout",
	}

	if viewer.UserID != 0 {
		if viewer.Token != nil {
			headerRenderer.CsrfToken = GenerateCSRFToken(viewer.Token.Token)
		}

		obj, err := f.App.Store.ObjGet(ctx, viewer.UserID)
		if err != nil {
			log.Printf("Error getting user: %s", err)
		} else if obj == nil || obj.User == nil {
			log.Printf("User %d not found", viewer.UserID)
		} else {
			user := obj.User

			headerRenderer.IsAuthorized = true
			headerRenderer.UserAvatar = user.VkAvatar
			headerRenderer.UserName = user.Name
		}
	}

	return &pb.FeedGetHeaderResponse{Renderer: &headerRenderer}, nil
}
