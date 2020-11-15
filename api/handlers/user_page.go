package handlers

import (
	"fmt"
	"strconv"

	"github.com/materkov/meme9/api/handlers/common"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/store"
)

type UserPage struct {
	Store *store.Store
}

func (p *UserPage) Handle(viewer *api.Viewer, req *pb.UserPageRequest) (*pb.UserPageRenderer, error) {
	userID, _ := strconv.Atoi(req.UserId)
	user, err := p.Store.GetUser(userID)
	if err == store.ErrNodeNotFound {
		return nil, api.NewError("USER_NOT_FOUND", "Пользователь не найден")
	} else if err != nil {
		return nil, fmt.Errorf("error getting user from store: %w", err)
	}

	renderer := &pb.UserPageRenderer{
		Id:             strconv.Itoa(user.ID),
		LastPostId:     "2",
		LastPostUrl:    "/posts/2",
		Name:           user.Name,
		HeaderRenderer: common.GetHeaderRenderer(viewer),
	}

	return renderer, nil
}
