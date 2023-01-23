package api

import (
	"context"
	"github.com/materkov/meme9/web5/imgproxy"
	"github.com/materkov/meme9/web5/pkg/files"
	"github.com/materkov/meme9/web5/store"
	"strconv"
	"strings"
)

type Photo struct {
	ID string `json:"id,omitempty"`

	Address string `json:"address,omitempty"`
	Width   int    `json:"width,omitempty"`
	Height  int    `json:"height,omitempty"`

	Thumbs []*PhotoThumb `json:"thumbs,omitempty"`
}

type PhotoThumb struct {
	Width  int `json:"width"`
	Height int `json:"height"`

	Address string `json:"address"`
}

func handlePhotosId(ctx context.Context, viewerID int, url string) []interface{} {
	photoID, _ := strconv.Atoi(strings.TrimPrefix(url, "/photos/"))

	result := Photo{
		ID: strconv.Itoa(photoID),
	}

	photo := store.CachedStoreFromCtx(ctx).Photo.Get(photoID)
	if photo == nil {
		return []interface{}{result}
	}

	result.Address = files.GetURL(photo.Hash)
	result.Width = photo.Width
	result.Height = photo.Height

	sizes := []int{50, 100, 300, 500, 1000}
	for _, size := range sizes {
		ratio := float64(size) / float64(photo.Width)
		if ratio > 1 {
			break
		}

		thumb := PhotoThumb{
			Width:   size,
			Height:  int(float64(photo.Height) * ratio),
			Address: imgproxy.GetURL(photo.Hash, size),
		}
		result.Thumbs = append(result.Thumbs, &thumb)
	}

	return []interface{}{result}
}
