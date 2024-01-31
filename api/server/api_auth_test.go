package server

import (
	"context"
	"fmt"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/pkg"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createTestDB(t *testing.T) func() {
	store2.GlobalStore = createMockStore()

	return func() {
	}
}

func createMockStore() *store2.Store {
	return &store2.Store{
		Unique: &store2.MockUniqueStore{
			Rows: map[string]int{},
		},
		Likes: &store2.MockLikes{
			Rows: map[string]bool{},
		},
		Subs: &store2.MockSubscriptions{
			Following: map[int][]int{},
		},
		Wall: &store2.MockWall{
			Posts: map[int][]int{},
		},
		Votes: &store2.MockVotes{
			Votes: map[int][]int{},
		},
		Users: &store2.MockUserStore{
			Objects: map[int]*store.User{},
		},
		Posts: &store2.MockPostStore{
			Objects: map[int]*store.Post{},
		},
		Polls: &store2.MockPollStore{
			Objects: map[int]*store.Poll{},
		},
		PollAnswers: &store2.MockPollAnswerStore{
			Objects: map[int]*store.PollAnswer{},
		},
		Tokens: &store2.MockTokenStore{
			Objects: map[int]*store.Token{},
		},
		Configs: &store2.MockConfigStore{
			Objects: map[int]*store.Config{},
		},
		Bookmarks: &store2.MockBookmarks{},
	}
}

func TestApi_authEmailLogin(t *testing.T) {
	srv := AuthServer{}
	closer := createTestDB(t)
	defer closer()

	ctx := context.Background()

	resp, err := srv.Register(ctx, &api.EmailReq{
		Email:    "test@email.com",
		Password: "12345",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)
	require.NotEmpty(t, resp.UserId)

	loginResp, err := srv.Login(ctx, &api.EmailReq{
		Email:    "test@email.com",
		Password: "12345",
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.Token)
	require.Equal(t, resp.UserId, loginResp.UserId)

	_, err = srv.Login(ctx, &api.EmailReq{
		Email:    "test@email.com",
		Password: "wrong password",
	})
	requireAPIError(t, err, "InvalidCredentials")

	_, err = srv.Login(ctx, &api.EmailReq{
		Email:    "bad@email.com",
		Password: "12345",
	})
	requireAPIError(t, err, "InvalidCredentials")

	_, err = srv.Login(ctx, &api.EmailReq{
		Email:    "",
		Password: "12345",
	})
	requireAPIError(t, err, "InvalidCredentials")

	_, err = srv.Login(ctx, &api.EmailReq{
		Email:    "test@email.com",
		Password: "",
	})
	requireAPIError(t, err, "InvalidCredentials")
}

func TestAPI_authRegister(t *testing.T) {
	srv := AuthServer{}
	closer := createTestDB(t)
	defer closer()

	resp, err := srv.Register(nil, &api.EmailReq{
		Email:    "test@mail.com",
		Password: "123",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)

	t.Run("email registered", func(t *testing.T) {
		_, err := srv.Register(nil, &api.EmailReq{Email: "test@mail.com", Password: "123"})
		requireAPIError(t, err, "EmailAlreadyRegistered")
	})
	t.Run("empty email", func(t *testing.T) {
		_, err := srv.Register(nil, &api.EmailReq{Email: ""})
		requireAPIError(t, err, "EmptyEmail")
	})
	t.Run("empty pass", func(t *testing.T) {
		_, err := srv.Register(nil, &api.EmailReq{Email: "test@email.com", Password: ""})
		requireAPIError(t, err, "EmptyPassword")
	})
	t.Run("email too long", func(t *testing.T) {
		_, err := srv.Register(nil, &api.EmailReq{Email: strings.Repeat("a", 1000), Password: "123"})
		requireAPIError(t, err, "EmailTooLong")
	})
}

func TestApi_authVK(t *testing.T) {
	srv := AuthServer{}
	closer := createTestDB(t)
	defer closer()

	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/access_token" {
			fmt.Fprint(w, `{"access_token": "vk_token", "user_id": 41512}`)
		} else {
			fmt.Fprint(w, `{"response": [{"id": 41512, "first_name": "Test", "last_name": "Testovich", "photo_200": "https://image.com/1.jpg"}]}`)
		}
	}))
	defer httpSrv.Close()

	pkg.VKEndpoint = httpSrv.URL
	pkg.VKAPIEndpoint = httpSrv.URL

	pkg.HTTPClient = httpSrv.Client()
	ctx := context.Background()

	resp1, err := srv.Vk(ctx, &api.VkReq{
		Code:        "1234",
		RedirectUrl: "https://site.me",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp1.Token)

	resp2, err := srv.Vk(ctx, &api.VkReq{
		Code:        "1234",
		RedirectUrl: "https://site.me",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp2.Token)
	require.Equal(t, resp1.UserId, resp2.UserId)
}
