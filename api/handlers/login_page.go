package handlers

import (
	"fmt"
	"log"
	"net/url"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers/common"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type LoginPage struct {
	Store *store.Store
}

func (l *LoginPage) Handle(viewer *api.Viewer, req *login.LoginPageRequest) *login.AnyRenderer {
	redirectURL := url.QueryEscape(fmt.Sprintf("%s://%s/vk-callback", viewer.RequestScheme, viewer.RequestHost))
	log.Printf("viewer %s", viewer.RequestHost)

	return &login.AnyRenderer{Renderer: &login.AnyRenderer_LoginPageRenderer{
		LoginPageRenderer: &login.LoginPageRenderer{
			SubmitUrl:      "/submit_url",
			WelcomeText:    "Login welcome текстик",
			HeaderRenderer: common.GetHeaderRenderer(viewer),
			VkUrl:          fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", VKAppID, redirectURL),
		},
	}}
}
