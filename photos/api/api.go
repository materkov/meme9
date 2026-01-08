package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/materkov/meme9/photos/processor"
)

type Processor interface {
	Process(ctx context.Context, file []byte) ([]byte, error)
}

type Uploader interface {
	Upload(ctx context.Context, file []byte) (url string, err error)
}

type API struct {
	processor Processor
	uploader  Uploader
}

func New(processor Processor, uploader Uploader) *API {
	return &API{processor: processor, uploader: uploader}
}

func (a *API) Start() error {
	mux := a.Routes()

	fmt.Println("Server is running on http://localhost:8081")
	return http.ListenAndServe("127.0.0.1:8081", mux)
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/twirp/meme.photos.Photos/upload", a.HandleUpload)

	return withCORS(mux)
}

func (a *API) HandleUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	resizedImg, err := a.processor.Process(ctx, bodyBytes)
	if errors.Is(err, processor.ErrInvalidImage) {
		http.Error(w, "Invalid image", http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("error processing image: %v", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	url, err := a.uploader.Upload(ctx, resizedImg)
	if err != nil {
		log.Printf("error uploading image: %v", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(url))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
