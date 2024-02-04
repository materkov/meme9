package api

import "github.com/materkov/meme9/web6/pb/github.com/materkov/meme9/api"

type Viewer struct {
	UserID       int
	UserName     string
	AuthToken    string
	IsCookieAuth bool
	ClientIP     string
}

var (
	ApiAuthClient  api.Auth
	ApiPostsClient api.Posts
	ApiPollsClient api.Polls
	ApiUsersClient api.Users

	ImageProxyClient api.ImageProxy
)
