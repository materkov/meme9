package handlers

import (
	"fmt"
	"log"
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
	postID, _ := strconv.Atoi(req.PostId)

	post, err := p.Store.GetPost(postID)
	if err == store.ErrNodeNotFound {
		return &login.AnyRenderer{Renderer: &login.AnyRenderer_ErrorRenderer{
			ErrorRenderer: &login.ErrorRenderer{
				ErrorCode:   "POST_NOT_FOUND",
				DisplayText: "Пост не найден",
			},
		}}
	} else if err != nil {
		log.Printf("[ERROR] Internal error: %s", err)
		return &login.AnyRenderer{Renderer: &login.AnyRenderer_ErrorRenderer{
			ErrorRenderer: &login.ErrorRenderer{
				ErrorCode:   "INTERNAL_ERROR",
				DisplayText: "Неизвестная ошибка",
			},
		}}
	}

	renderer := &login.PostPageRenderer{
		Id:             req.PostId,
		PostUrl:        fmt.Sprintf("/posts/%s", req.PostId),
		Text:           post.Text,
		UserId:         strconv.Itoa(post.UserID),
		UserUrl:        fmt.Sprintf("/users/%d", post.UserID),
		HeaderRenderer: common.GetHeaderRenderer(viewer),
	}

	return &login.AnyRenderer{Renderer: &login.AnyRenderer_PostPageRenderer{
		PostPageRenderer: renderer,
	}}
}
