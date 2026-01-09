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
type Auth interface {
	Auth(ctx context.Context, header string) (string, error)
}

type Service struct {
	processor Processor
	uploader  Uploader
	auth      Auth
}

func New(processor Processor, uploader Uploader, auth Auth) *Service {
	return &Service{
		processor: processor,
		uploader:  uploader,
		auth:      auth,
	}
}

func (s *Service) Start() error {
	mux := s.Routes()

	fmt.Println("Server is running on http://localhost:8081")
	return http.ListenAndServe("127.0.0.1:8081", mux)
}

func (s *Service) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/twirp/meme.photos.Photos/upload", s.HandleUpload)

	return withCORS(mux)
}

func CORSMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (s *Service) HandleUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := s.auth.Auth(ctx, r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "auth_required", http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	_ = r.Body.Close()

	resizedImg, err := s.processor.Process(ctx, bodyBytes)
	if errors.Is(err, processor.ErrInvalidImage) {
		http.Error(w, "Invalid image", http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("error processing image: %v", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	url, err := s.uploader.Upload(ctx, resizedImg)
	if err != nil {
		log.Printf("error uploading image: %v", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(url))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
