package handlers

import (
	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers/common"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type Composer struct {
	Store *store.Store
}

func (c *Composer) Handle(viewer *api.Viewer, req *pb.ComposerRequest) (*pb.ComposerRenderer, error) {
	renderer := &pb.ComposerRenderer{
		HeaderRenderer: common.GetHeaderRenderer(viewer),
		WelcomeText:    "текст фром бакеэнд",
	}

	if viewer.User == nil {
		renderer.UnathorizedText = "Нужно авторизоваться чтобы написать пост"
	}

	return renderer, nil
}
