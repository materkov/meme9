package handlers

import (
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/pkg/config"
	"github.com/materkov/meme9/api/store"
)

const VKAppID = 7260220

type Handlers struct {
	store  *store.Store
	Config *config.Config

	loginPage *LoginPage
	addPost   *AddPost
	getFeed   *GetFeed
	composer  *Composer
	index     *Index
	postPage  *PostPage
	userPage  *UserPage
}

func NewHandlers(store *store.Store, config *config.Config) *Handlers {
	h := &Handlers{store: store, Config: config}

	h.loginPage = &LoginPage{Store: store}
	h.addPost = &AddPost{Store: store}
	h.postPage = &PostPage{Store: store}
	h.userPage = &UserPage{Store: store}
	h.getFeed = &GetFeed{Store: store}
	h.composer = &Composer{Store: store}
	h.index = &Index{Store: store}

	return h
}

func (h *Handlers) Call(viewer *api.Viewer, method string, args string) (proto.Message, error) {
	switch method {
	case "meme.API.UserPage":
		req := &pb.UserPageRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return nil, api.NewError("INVALID_REQUEST", "Failed parsing request")
		}
		return h.userPage.Handle(viewer, req)
	case "meme.API.PostPage":
		req := &pb.PostPageRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return nil, api.NewError("INVALID_REQUEST", "Failed parsing request")
		}
		return h.postPage.Handle(viewer, req)
	case "meme.API.LoginPage":
		req := &pb.LoginPageRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return nil, api.NewError("INVALID_REQUEST", "Failed parsing request")
		}
		return h.loginPage.Handle(viewer, req)
	case "meme.API.AddPost":
		req := &pb.AddPostRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return nil, api.NewError("INVALID_REQUEST", "Failed parsing request")
		}
		return h.addPost.Handle(viewer, req)
	case "meme.API.GetFeed":
		req := &pb.GetFeedRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return nil, api.NewError("INVALID_REQUEST", "Failed parsing request")
		}
		return h.getFeed.Handle(viewer, req)
	case "meme.API.Composer":
		req := &pb.ComposerRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return nil, api.NewError("INVALID_REQUEST", "Failed parsing request")
		}
		return h.composer.Handle(viewer, req)
	case "meme.API.Index":
		req := &pb.IndexRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return nil, api.NewError("INVALID_REQUEST", "Failed parsing request")
		}
		return h.index.Handle(viewer, req)
	default:
		return nil, api.NewError("INVALID_METHOD", "Method not found")
	}
}
