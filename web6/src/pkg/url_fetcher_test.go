package pkg

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
<html><head>
	<meta property="og:title" content="Test title" />
	<meta property="og:description" content="Test description" />
	<meta property="og:image" content="https://image.com/example.jpg" />
</head></html>
`)
	}))
	defer srv.Close()

	data, err := FetchURL(srv.URL)
	require.NoError(t, err)
	require.Equal(t, "Test title", data.Title)
	require.Equal(t, "Test description", data.Description)
	require.Equal(t, "https://image.com/example.jpg", data.ImageURL)
	require.Equal(t, srv.URL, data.FinalURL)
	require.Equal(t, srv.URL, data.URL)
}
