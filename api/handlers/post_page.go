package handlers

import (
	"fmt"
	"strconv"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers/common"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type PostPage struct {
	Store *store.Store
}

func (p *PostPage) Handle(viewer *api.Viewer, req *pb.PostPageRequest) (*pb.PostPageRenderer, error) {
	postID, _ := strconv.Atoi(req.PostId)

	post, err := p.Store.GetPost(postID)
	if err == store.ErrNodeNotFound {
		return nil, api.NewError("POST_NOT_FOUND", "Пост не найден")
	} else if err != nil {
		return nil, fmt.Errorf("error getting post from store: %w", err)
	}

	renderer := &pb.PostPageRenderer{
		Id:             req.PostId,
		PostUrl:        fmt.Sprintf("/posts/%s", req.PostId),
		Text:           post.Text,
		UserId:         strconv.Itoa(post.UserID),
		UserUrl:        fmt.Sprintf("/users/%d", post.UserID),
		HeaderRenderer: common.GetHeaderRenderer(viewer),
	}

	return renderer, nil
}
