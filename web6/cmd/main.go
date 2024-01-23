package main

import (
	"context"
	"fmt"
	pbapi "github.com/materkov/meme9/web6/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/web6/src/api"
	"github.com/twitchtv/twirp"
	"net/http"
)

func main() {
	hooks := twirp.ClientHooks{
		RequestPrepared: func(ctx context.Context, request *http.Request) (context.Context, error) {
			viewer, ok := ctx.Value(api.CtxViewer).(*api.Viewer)
			if ok && viewer.IsCookieAuth {
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", viewer.AuthToken))
			}
			return ctx, nil
		},
	}

	api.ApiAuthClient = pbapi.NewAuthProtobufClient("http://localhost:8002", http.DefaultClient, twirp.WithClientHooks(&hooks))
	api.ApiUsersClient = pbapi.NewUsersProtobufClient("http://localhost:8002", http.DefaultClient, twirp.WithClientHooks(&hooks))
	api.ApiPollsClient = pbapi.NewPollsProtobufClient("http://localhost:8002", http.DefaultClient, twirp.WithClientHooks(&hooks))
	api.ApiPostsClient = pbapi.NewPostsProtobufClient("http://localhost:8002", http.DefaultClient, twirp.WithClientHooks(&hooks))

	s := &api.HttpServer{}
	s.Serve()
}
