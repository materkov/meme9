package handlers

import (
	"github.com/materkov/meme9/api/api"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type Composer struct {
	Store *store.Store
}

func (c *Composer) Handle(viewer *api.Viewer, req *login.ComposerRequest) *login.AnyRenderer {
	return &login.AnyRenderer{Renderer: &login.AnyRenderer_ComposerRenderer{
		ComposerRenderer: &login.ComposerRenderer{
			WelcomeText: "текст фром бакеэнд",
		},
	}}
}
