package handlers

import (
	"github.com/materkov/meme9/api/api"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type LoginPage struct {
	Store *store.Store
}

func (l *LoginPage) Handle(viewer *api.Viewer, req *login.LoginPageRequest) *login.AnyRenderer {
	return &login.AnyRenderer{Renderer: &login.AnyRenderer_LoginPageRenderer{
		LoginPageRenderer: &login.LoginPageRenderer{
			SubmitUrl:   "/submit_url",
			WelcomeText: "Login welcome текстик",
		},
	}}
}
