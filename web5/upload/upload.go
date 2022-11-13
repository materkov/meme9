package upload

import (
	"bytes"
	"encoding/json"
	"github.com/materkov/meme9/web5/pkg/files"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"net/http"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")

	if r.Method == "OPTIONS" {
		return
	}

	fileBytes, err := io.ReadAll(r.Body)
	if err != nil || len(fileBytes) == 0 {
		w.WriteHeader(400)
		return
	}

	// 3 MBytes
	if len(fileBytes) > 3*1024*1024 {
		w.WriteHeader(400)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width < 10 || width > 10000 || height < 10 || height > 10000 {
		w.WriteHeader(400)
		return
	}

	fileHash := files.Hash(fileBytes)

	err = files.SelectelUpload(fileBytes, fileHash)
	if err != nil {
		log.Printf("[ERROR] Error uplaoding file: %s", err)
		w.WriteHeader(400)
		return
	}

	resp := struct {
		UploadToken string `json:"uploadToken"`
	}{
		UploadToken: fileHash,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
