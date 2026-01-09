package api

import (
	"context"

	"github.com/twitchtv/twirp"

	proto "github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api/photos"
)

// Photos implements the photos Twirp service.
// Uploads are handled by the dedicated photos service endpoint.
type Photos struct{}

// NewPhotos constructs a Photos service instance.
func NewPhotos() *Photos {
	return &Photos{}
}

// GenerateUploadUrl is disabled.
//
// Uploads should go directly to the photos service:
// POST https://meme2.mmaks.me/twirp/meme.photos.Photos/upload
// The resulting public URL is returned in the response body.
func (s *Photos) GenerateUploadUrl(ctx context.Context, req *proto.GenerateUploadUrlRequest) (*proto.GenerateUploadUrlResponse, error) {
	return nil, twirp.NewError(twirp.Unimplemented, "GenerateUploadUrl is disabled; upload via photos service endpoint")
}
