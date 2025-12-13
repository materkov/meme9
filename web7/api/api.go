package api

import (
	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/subscriptions"
	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	postsservice "github.com/materkov/meme9/web7/services/posts"
	tokensservice "github.com/materkov/meme9/web7/services/tokens"
)

type API struct {
	posts         *posts.Adapter
	users         *users.Adapter
	tokens        *tokens.Adapter
	subscriptions *subscriptions.Adapter

	postsService  *postsservice.Service
	tokensService *tokensservice.Service
}

func NewAPI(postsAdapter *posts.Adapter, usersAdapter *users.Adapter, tokensAdapter *tokens.Adapter, subscriptionsAdapter *subscriptions.Adapter, postsService *postsservice.Service, tokensService *tokensservice.Service) *API {
	return &API{
		posts:         postsAdapter,
		users:         usersAdapter,
		tokens:        tokensAdapter,
		subscriptions: subscriptionsAdapter,
		postsService:  postsService,
		tokensService: tokensService,
	}
}
