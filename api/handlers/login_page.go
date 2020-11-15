package handlers

import (
	"fmt"
	"net/url"

	"github.com/materkov/meme9/api/handlers/common"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/store"
)

type LoginPage struct {
	Store *store.Store
}

func (l *LoginPage) Handle(viewer *api.Viewer, req *pb.LoginPageRequest) (*pb.LoginPageRenderer, error) {
	redirectURL := url.QueryEscape(fmt.Sprintf("%s://%s/vk-callback", viewer.RequestScheme, viewer.RequestHost))

	renderer := &pb.LoginPageRenderer{
		SubmitUrl:      "/submit_url",
		WelcomeText:    "Login welcome текстик",
		HeaderRenderer: common.GetHeaderRenderer(viewer),
		VkUrl:          fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", VKAppID, redirectURL),
		VkText:         "Войти через ВК",
	}

	return renderer, nil
}
