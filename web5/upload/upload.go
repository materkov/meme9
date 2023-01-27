package upload

import (
	"bytes"
	"encoding/json"
	"github.com/materkov/meme9/web5/pkg/files"
	"github.com/materkov/meme9/web5/store"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")

	if r.Method == "OPTIONS" {
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
		return
	}

	fileBytes, err := io.ReadAll(file)
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

	photo := store.Photo{
		ID:     int(time.Now().Unix()),
		Size:   len(fileBytes),
		Hash:   fileHash,
		Width:  width,
		Height: height,
	}

	err = store.NodeSave(photo.ID, store.ObjectTypePhoto, &photo)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	resp := struct {
		UploadToken string `json:"uploadToken"`
	}{
		UploadToken: strconv.Itoa(photo.ID),
	}
	_ = json.NewEncoder(w).Encode(resp)
}
