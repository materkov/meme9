package handlers

import (
	"strconv"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers/common"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type PostPage struct {
	Store *store.Store
}

func (p *PostPage) Handle(viewer *api.Viewer, req *login.PostPageRequest) *login.AnyRenderer {
	return &login.AnyRenderer{Renderer: &login.AnyRenderer_PostPageRenderer{
		PostPageRenderer: &login.PostPageRenderer{
			Id:             req.PostId,
			Text:           "bla bla bla - " + req.PostId,
			UserId:         "1",
			CurrentUserId:  strconv.Itoa(viewer.UserID),
			HeaderRenderer: common.GetHeaderRenderer(viewer),
		},
	}}
}
