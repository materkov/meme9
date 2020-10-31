package handlers

import (
	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers/common"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type Index struct {
	Store *store.Store
}

func (i *Index) Handle(viewer *api.Viewer, req *login.IndexRequest) *login.AnyRenderer {
	return &login.AnyRenderer{Renderer: &login.AnyRenderer_IndexRenderer{
		IndexRenderer: &login.IndexRenderer{
			Text:           "Текст главной странцы сайтика",
			FeedUrl:        "/feed",
			HeaderRenderer: common.GetHeaderRenderer(viewer),
		},
	}}
}
