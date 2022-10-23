package app

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
)

import (
	"github.com/stretchr/testify/require"
)

func TestHandleVkCallback(t *testing.T) {
	setupDB(t)
	app := App{Store: ObjectStore}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST", "https://oauth.vk.com/access_token",
		httpmock.NewStringResponder(200, `{"access_token": "vk-test-token", "user_id": 55161}`),
	)
	httpmock.RegisterResponder(
		"POST", "https://api.vk.com/method/users.get",
		httpmock.NewStringResponder(200, `{"response":[{"first_name": "Maks", "last_name": "Materkov", "photo_200": "https://test.com/image1"}]}`),
	)

	token, err := app.DoVKCallback(context.Background(), "test-code", &Viewer{})
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestHandleVkCallback_VkError(t *testing.T) {
	app := App{Store: ObjectStore}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://oauth.vk.com/access_token",
		httpmock.NewStringResponder(200, `{"error": "error"}`),
	)

	_, err := app.DoVKCallback(context.Background(), "test-code", &Viewer{})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "empty access token")
}