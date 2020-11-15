package handlers

import (
	"github.com/materkov/meme9/api/handlers/common"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/store"
)

type Index struct {
	Store *store.Store
}

func (i *Index) Handle(viewer *api.Viewer, req *pb.IndexRequest) (*pb.IndexRenderer, error) {
	renderer := &pb.IndexRenderer{
		Text:           "Текст главной странцы сайтика",
		FeedUrl:        "/feed",
		HeaderRenderer: common.GetHeaderRenderer(viewer),
	}

	return renderer, nil
}
