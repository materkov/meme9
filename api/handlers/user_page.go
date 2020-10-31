package handlers

import (
	"strconv"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers/common"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type UserPage struct {
	Store *store.Store
}

func (p *UserPage) Handle(viewer *api.Viewer, req *login.UserPageRequest) *login.AnyRenderer {
	return &login.AnyRenderer{Renderer: &login.AnyRenderer_UserPageRenderer{
		UserPageRenderer: &login.UserPageRenderer{
			Id:             req.UserId,
			LastPostId:     "2",
			LastPostUrl:    "/posts/2",
			CurrentUserId:  strconv.Itoa(viewer.UserID),
			Name:           req.UserId + " - name",
			HeaderRenderer: common.GetHeaderRenderer(viewer),
		},
	}}
}
