package imgproxy

import (
	"fmt"
	"github.com/materkov/meme9/web5/pkg/files"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")

	hash := r.URL.Query().Get("hash")
	if hash == "" {
		log.Printf("Empty hash")
		return
	}

	width, _ := strconv.Atoi(r.URL.Query().Get("w"))
	if width <= 0 {
		log.Printf("Empty width")
		return
	}

	imgUrl := files.GetURL(hash)

	resp, err := http.Get(imgUrl)
	if err != nil {
		log.Printf("Error loading image: %s", err)
		return
	} else if resp.StatusCode != 200 {
		log.Printf("Bad HTTP code: %d", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	src, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("Not an image: %s", err)
		return
	}

	originalWidth := src.Bounds().Max.X
	originalHeight := src.Bounds().Max.Y

	ratio := float64(width) / float64(originalWidth)
	if ratio > 1 {
		ratio = 1
	}

	resizedWidth := int(float64(originalWidth) * ratio)
	resizedHeight := int(float64(originalHeight) * ratio)

	dst := image.NewRGBA(image.Rect(0, 0, resizedWidth, resizedHeight))
	draw.ApproxBiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	_ = jpeg.Encode(w, dst, nil)
}

func GetURL(hash string, width int) string {
	return fmt.Sprintf("https://meme.mmaks.me/imgproxy?hash=%s&w=%d", hash, width)
}
