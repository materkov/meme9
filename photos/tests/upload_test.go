package tests

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/materkov/meme9/photos/api"
	"github.com/materkov/meme9/photos/processor"
	"github.com/materkov/meme9/photos/uploader"
	"github.com/stretchr/testify/require"
)

const textImage = "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD3+iiigD//2Q=="

func setupAPI(t *testing.T) *api.API {
	godotenv.Load()

	processor := processor.New()
	uploader, err := uploader.New(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_KEY"),
	)
	require.NoError(t, err)
	api := api.New(processor, uploader)

	return api
}

func TestUploadFile(t *testing.T) {
	api := setupAPI(t)
	testSrv := httptest.NewServer(api.Routes())
	defer testSrv.Close()

	textImageBytes, err := base64.StdEncoding.DecodeString(textImage)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testSrv.URL+"/twirp/meme.photos.Photos/upload", bytes.NewReader(textImageBytes))
	require.NoError(t, err)

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	uploadedURL := strings.TrimSpace(string(respBytes))
	require.NotEmpty(t, uploadedURL)

	downloadReq, err := http.NewRequest(http.MethodGet, uploadedURL, nil)
	require.NoError(t, err)

	downloadResp, err := client.Do(downloadReq)
	require.NoError(t, err)
	defer downloadResp.Body.Close()

	require.Equal(t, http.StatusOK, downloadResp.StatusCode)

	downloadedBytes, err := io.ReadAll(downloadResp.Body)
	require.NoError(t, err)
	require.NotEmpty(t, downloadedBytes)
	log.Printf("Downloaded bytes from: %s", uploadedURL)
}
