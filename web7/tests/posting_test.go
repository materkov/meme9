package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
	"github.com/materkov/meme9/web7/api"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const baseURL = "http://localhost:8080"

var testAPI *api.API

func initAPI() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27017/meme9?authSource=admin"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize adapters
	postsAdapter := posts.NewAdapter(client)
	usersAdapter := users.NewAdapter(client)
	tokensAdapter := tokens.NewAdapter(client)

	testAPI = api.NewAPI(postsAdapter, usersAdapter, tokensAdapter)
}

func apiRequest(t *testing.T, method string, req, resp any) {
	body, err := json.Marshal(req)
	require.NoError(t, err)
	respBody, err := http.Post(baseURL+method, "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer respBody.Body.Close()

	err = json.NewDecoder(respBody.Body).Decode(resp)
	require.NoError(t, err)
}

func TestPosting(t *testing.T) {
	if testAPI == nil {
		initAPI()
	}
	go testAPI.Serve()

	req := api.PublishReq{
		Text: "test text",
	}
	resp := api.PublishResp{}
	apiRequest(t, "/publish", &req, &resp)
	require.Empty(t, resp.ID)
}
