package api

import (
	"context"
	"github.com/materkov/meme9/web5/pkg/files"
	"github.com/materkov/meme9/web5/store"
	"strconv"
	"strings"
)

type Photo struct {
	URL string `json:"url,omitempty"`
	ID  string `json:"id,omitempty"`

	Address string `json:"address,omitempty"`
	Width   int    `json:"width,omitempty"`
	Height  int    `json:"height,omitempty"`
}

func handlePhotosId(ctx context.Context, viewerID int, url string) []interface{} {
	photoID, _ := strconv.Atoi(strings.TrimPrefix(url, "/photos/"))

	result := Photo{
		URL: url,
		ID:  strconv.Itoa(photoID),
	}

	photo := store.CachedStoreFromCtx(ctx).Photo.Get(photoID)
	if photo == nil {
		return []interface{}{result}
	}

	result.Address = files.GetURL(photo.Hash)
	result.Width = photo.Width
	result.Height = photo.Height

	return []interface{}{result}
}
