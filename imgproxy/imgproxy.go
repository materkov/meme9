package imgproxy

import (
	"bytes"
	"context"
	"fmt"
	"github.com/materkov/meme9/imgproxy/pb/github.com/materkov/meme9/api"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"net/http"
)

type Service struct{}

func (s *Service) Resize(ctx context.Context, req *api.ResizeReq) (*api.ResizeResp, error) {
	resp, err := http.Get(req.ImageUrl)
	if err != nil {
		return nil, fmt.Errorf("error doing http request: %w", err)
	}

	src, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error decoding jpeg: %w", err)
	}

	_ = resp.Body.Close()
	result := bytes.NewBuffer(nil)

	if src.Bounds().Size().X <= 200 {
		err := jpeg.Encode(result, src, nil)
		if err != nil {
			return nil, fmt.Errorf("error encoding result jpeg: %w", err)
		}
	} else {
		dst := image.NewRGBA(image.Rect(0, 0, 200, 200))
		draw.BiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

		err = jpeg.Encode(result, dst, nil)
		if err != nil {
			return nil, fmt.Errorf("error encoding resized jpeg: %w", err)
		}
	}

	return &api.ResizeResp{Image: result.Bytes()}, nil
}
