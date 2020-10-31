package common

import (
	"strconv"

	"github.com/materkov/meme9/api/api"
	login "github.com/materkov/meme9/api/pb"
)

func GetHeaderRenderer(viewer *api.Viewer) *login.HeaderRenderer {
	renderer := login.HeaderRenderer{}

	if viewer.UserID != 0 {
		renderer.CurrentUserId = strconv.Itoa(viewer.User.ID)
		renderer.CurrentUserName = viewer.User.Name
	}

	renderer.Links = []*login.HeaderRenderer_Link{
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
