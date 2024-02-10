package server

import (
	"context"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/pkg"
	"github.com/twitchtv/twirp"
)

type PhotosServer struct{}

func (p *PhotosServer) Upload(_ context.Context, req *api.UploadReq) (*api.UploadResp, error) {
	width, height, err := pkg.ValidatePhoto(req.PhotoBytes)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "InvalidImage")
	}

	fileHash := pkg.GetFileHash(req.PhotoBytes)
	err = pkg.SelectelUpload(req.PhotoBytes, fileHash)
	if err != nil {
		return nil, err
	}

	token := pkg.UploadToken{
		Hash:   fileHash,
		Width:  width,
		Height: height,
		Size:   len(req.PhotoBytes),
	}

	return &api.UploadResp{UploadToken: token.ToString()}, nil
}
