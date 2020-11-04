package handlers

import (
	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers/common"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type Composer struct {
	Store *store.Store
}

func (c *Composer) Handle(viewer *api.Viewer, req *login.ComposerRequest) *login.AnyRenderer {
	renderer := &login.ComposerRenderer{
		HeaderRenderer: common.GetHeaderRenderer(viewer),
	}

	if viewer.User == nil {
		renderer.UnathorizedText = "Нужно авторизоваться чтобы написать пост"
	} else {
		renderer.WelcomeText = "текст фром бакеэнд"
	}

	return &login.AnyRenderer{Renderer: &login.AnyRenderer_ComposerRenderer{
		ComposerRenderer: renderer,
	}}
}
