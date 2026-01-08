//go:build integration
// +build integration

package tests

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/materkov/meme9/photos/api"
	"github.com/materkov/meme9/photos/auth"
	"github.com/materkov/meme9/photos/processor"
	"github.com/materkov/meme9/photos/uploader"
	"github.com/stretchr/testify/require"

	pb "github.com/materkov/meme9/photos/internal/authclient/pb/github.com/materkov/meme9/api/auth"
)

const textImage = "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD3+iiigD//2Q=="

func setupAPI(t *testing.T) (string, func()) {
	godotenv.Load()

	processor := processor.New()
	uploader, err := uploader.New(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_KEY"),
	)
	require.NoError(t, err)

	mockSrv := httptest.NewServer(pb.NewAuthServer(&MockAuth{}))
	authClient := pb.NewAuthJSONClient(mockSrv.URL, http.DefaultClient)
	authService := auth.New(authClient)

	api := api.New(processor, uploader, authService)
	testSrv := httptest.NewServer(api.Routes())

	return testSrv.URL, func() {
		mockSrv.Close()
		testSrv.Close()
	}
}

func requireImageExists(t *testing.T, url string) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	downloadedBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NotEmpty(t, downloadedBytes)
}

func TestUploadFile(t *testing.T) {
	textImageBytes, err := base64.StdEncoding.DecodeString(textImage)
	require.NoError(t, err)

	testSrv, cleanup := setupAPI(t)
	defer cleanup()

	req, err := http.NewRequest(http.MethodPost, testSrv+"/twirp/meme.photos.Photos/upload", bytes.NewReader(textImageBytes))
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer good-token")

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	requireImageExists(t, string(respBytes))
}
