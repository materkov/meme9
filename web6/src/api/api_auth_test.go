package api

import (
	"database/sql"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
	_ "github.com/materkov/meme9/web6/src/store/sqlmock"
	"github.com/materkov/meme9/web6/src/store2"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func createTestDB(t *testing.T) func() {
	db, err := sql.Open("sqlmock", "TestDB_"+strconv.Itoa(rand.Int()))
	require.NoError(t, err)

	store.GlobalStore = &store.SqlStore{DB: db}
	store2.GlobalStore = createMockStore()

	return func() {
		_ = db.Close()
	}
}

func createMockStore() *store2.Store {
	return &store2.Store{
		Unique: store2.NewMockUniqueStore(),
	}
}

func createAPI() *API {
	return &API{}
}

func TestApi_authEmailLogin(t *testing.T) {
	api := createAPI()
	closer := createTestDB(t)
	defer closer()

	resp, err := api.authRegister(&Viewer{}, &AuthEmailReq{
		Email:    "test@email.com",
		Password: "12345",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)
	require.NotEmpty(t, resp.UserID)

	loginResp, err := api.authLogin(&Viewer{}, &AuthEmailReq{
		Email:    "test@email.com",
		Password: "12345",
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.Token)
	require.Equal(t, resp.UserID, loginResp.UserID)

	_, err = api.authLogin(&Viewer{}, &AuthEmailReq{
		Email:    "test@email.com",
		Password: "wrong password",
	})
	requireAPIError(t, err, "InvalidCredentials")

	_, err = api.authLogin(&Viewer{}, &AuthEmailReq{
		Email:    "bad@email.com",
		Password: "12345",
	})
	requireAPIError(t, err, "InvalidCredentials")

	_, err = api.authLogin(&Viewer{}, &AuthEmailReq{
		Email:    "",
		Password: "12345",
	})
	requireAPIError(t, err, "InvalidCredentials")

	_, err = api.authLogin(&Viewer{}, &AuthEmailReq{
		Email:    "test@email.com",
		Password: "",
	})
	requireAPIError(t, err, "InvalidCredentials")
}

func TestAPI_authRegister(t *testing.T) {
	api := createAPI()
	closer := createTestDB(t)
	defer closer()

	resp, err := api.authRegister(nil, &AuthEmailReq{
		Email:    "test@mail.com",
		Password: "123",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)

	t.Run("email registered", func(t *testing.T) {
		_, err := api.authRegister(nil, &AuthEmailReq{Email: "test@mail.com", Password: "123"})
		requireAPIError(t, err, "EmailAlreadyRegistered")
	})
	t.Run("empty email", func(t *testing.T) {
		_, err := api.authRegister(nil, &AuthEmailReq{Email: ""})
		requireAPIError(t, err, "EmptyEmail")
	})
	t.Run("empty pass", func(t *testing.T) {
		_, err := api.authRegister(nil, &AuthEmailReq{Email: "test@email.com", Password: ""})
		requireAPIError(t, err, "EmptyPassword")
	})
	t.Run("email too long", func(t *testing.T) {
		_, err := api.authRegister(nil, &AuthEmailReq{Email: strings.Repeat("a", 1000), Password: "123"})
		requireAPIError(t, err, "EmailTooLong")
	})
}

func TestApi_authVK(t *testing.T) {
	api := createAPI()
	closer := createTestDB(t)
	defer closer()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/access_token" {
			fmt.Fprint(w, `{"access_token": "vk_token", "user_id": 41512}`)
		} else {
			fmt.Fprint(w, `{"response": [{"id": 41512, "first_name": "Test", "last_name": "Testovich", "photo_200": "https://image.com/1.jpg"}]}`)
		}
	}))
	defer srv.Close()

	pkg.VKEndpoint = srv.URL
	pkg.VKAPIEndpoint = srv.URL

	pkg.HTTPClient = srv.Client()

	resp1, err := api.authVk(&Viewer{}, &AuthVkReq{
		Code:        "1234",
		RedirectURL: "https://site.me",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp1.Token)

	resp2, err := api.authVk(&Viewer{}, &AuthVkReq{
		Code:        "1234",
		RedirectURL: "https://site.me",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp2.Token)
	require.Equal(t, resp1.UserID, resp2.UserID)
}
