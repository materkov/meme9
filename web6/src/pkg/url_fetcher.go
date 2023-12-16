package pkg

import (
	"bytes"
	"fmt"
	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"io"
	"net/http"
	"regexp"
)

func FetchURL(url string) (*store.PostLink, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "curl/7.54")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	og := opengraph.NewOpenGraph()
	err = og.ProcessHTML(bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error processing html: %w", err)
	}

	link := store.PostLink{
		URL:         url,
		Title:       og.Title,
		Description: og.Description,
		FinalURL:    resp.Request.URL.String(),
	}

	if len(og.Images) > 0 {
		link.ImageURL = og.Images[0].URL
	}

	return &link, nil
}

func TryParseLink(post *store.Post) error {
	re := regexp.MustCompile(`https?://[^ ]+`)
	foundUrl := re.FindString(post.Text)
	if foundUrl == "" {
		return nil
	}

	linkInfo, err := FetchURL(foundUrl)
	if err != nil {
		return err
	}

	post.Link = linkInfo
	err = store2.GlobalStore.Posts.Update(post)
	if err != nil {
		return err
	}

	return nil
}
