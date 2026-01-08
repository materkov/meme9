package processor

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"golang.org/x/image/draw"
)

type Processor struct {
}

func New() *Processor {
	return &Processor{}
}

var ErrInvalidImage = fmt.Errorf("invalid image")

func (p *Processor) Process(ctx context.Context, file []byte) ([]byte, error) {
	contentType := http.DetectContentType(file)

	var img image.Image
	var err error

	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(file))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(file))
	default:
		err = fmt.Errorf("unsupperted content-type")
	}
	if err != nil {
		return nil, ErrInvalidImage
	}

	imgResized := resizeDownMax(img)

	var out bytes.Buffer
	err = jpeg.Encode(&out, imgResized, &jpeg.Options{Quality: 95})
	if err != nil {
		return nil, fmt.Errorf("failed to encode resized image: %w", err)
	}

	return out.Bytes(), nil
}

func resizeDownMax(src image.Image) image.Image {
	const maxW = 400
	const maxH = 400

	b := src.Bounds()
	w := b.Dx()
	h := b.Dy()

	if w <= maxW && h <= maxH {
		return src
	}

	scaleW := float64(maxW) / float64(w)
	scaleH := float64(maxH) / float64(h)
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}

	newW := int(float64(w) * scale)
	newH := int(float64(h) * scale)
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Over, nil)
	return dst
}
