package pkg

import (
	"fmt"
	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/materkov/meme9/web6/src/store"
	"net/http"
	"regexp"
)

func FetchURL(url string) (*store.PostLink, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	og := opengraph.NewOpenGraph()
	err = og.ProcessHTML(resp.Body)
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
	re := regexp.MustCompile("(http|https)://([\\w_-]+(?:(?:\\.[\\w_-]+)+))([\\w.,@?^=%&:/~+#-]*[\\w@?^=%&/~+#-])")
	foundUrl := re.FindString(post.Text)
	if foundUrl == "" {
		return nil
	}

	linkInfo, err := FetchURL(foundUrl)
	if err != nil {
		return err
	}

	post.Link = linkInfo
	err = store.UpdateObject(post, post.ID)
	if err != nil {
		return err
	}

	return nil
}
