package server

import (
	"context"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/pkg"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/twitchtv/twirp"
	"strconv"
)

func transformFile(file *store.File) *api.File {
	return &api.File{
		Url:    pkg.GetFilePath(file.Hash),
		Width:  int32(file.PhotoWidth),
		Height: int32(file.PhotoHeight),
	}
}

type PhotosServer struct{}

func (p *PhotosServer) Upload(ctx context.Context, req *api.UploadReq) (*api.UploadResp, error) {
	width, height, err := pkg.ValidatePhoto(req.PhotoBytes)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "InvalidImage")
	}

	viewer := ctx.Value(CtxViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "NotAuthorized")
	}

	fileHash := pkg.GetFileHash(req.PhotoBytes)
	err = pkg.SelectelUpload(req.PhotoBytes, fileHash)
	if err != nil {
		return nil, err
	}

	file := store.File{
		Hash:        fileHash,
		PhotoWidth:  width,
		PhotoHeight: height,
		Size:        len(req.PhotoBytes),
		UserID:      viewer.UserID,
	}

	err = store2.GlobalStore.Files.Add(&file)
	if err != nil {
		return nil, err
	}

	return &api.UploadResp{UploadToken: strconv.Itoa(file.ID)}, nil
}
