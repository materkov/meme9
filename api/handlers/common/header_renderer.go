package common

import (
	"strconv"

	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
)

func GetHeaderRenderer(viewer *api.Viewer) *pb.HeaderRenderer {
	renderer := pb.HeaderRenderer{}

	if viewer.User != nil {
		renderer.CurrentUserId = strconv.Itoa(viewer.User.ID)
		renderer.CurrentUserName = viewer.User.Name
	}

	renderer.Links = []*pb.HeaderRenderer_Link{
		{
			Url:   "/",
			Label: "Главная страница",
		},
		{
			Url:   "/feed",
			Label: "Лента",
		},
		{
			Url:   "/login",
			Label: "Логин",
		},
		{
			Url:   "/composer",
			Label: "Написать пост",
		},
	}

	return &renderer
}
