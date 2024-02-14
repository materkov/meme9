package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/twitchtv/twirp"
	"io"
	"log"
	"net/http"
	"strings"
)

func (h *HttpServer) ApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Version", pkg.BuildTime)

	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}

	path := strings.Replace(r.URL.Path, "/api/", "/twirp/", 1)

	apiReq, _ := http.NewRequest("POST", "http://localhost:8002"+path, r.Body)
	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Authorization", r.Header.Get("Authorization"))

	resp, err := http.DefaultClient.Do(apiReq)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// TODO think about error format
		respBody := struct {
			Code  string `json:"code,omitempty"`
			Msg   string `json:"msg,omitempty"`
			Error string `json:"error,omitempty"`
		}{}
		_ = json.NewDecoder(resp.Body).Decode(&respBody)

		respBody.Code = ""
		respBody.Error = respBody.Msg
		respBody.Msg = ""

		_ = json.NewEncoder(w).Encode(respBody)
	} else {
		_, _ = io.Copy(w, resp.Body)
	}
}

func (h *HttpServer) UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
		return
	}

	fileBytes, _ := io.ReadAll(file)

	ctx, _ := twirp.WithHTTPRequestHeaders(context.Background(), http.Header{
		"authorization": r.Header.Values("authorization"),
	})

	resp, err := ApiPhotosClient.Upload(ctx, &api.UploadReq{PhotoBytes: fileBytes})
	if err != nil {
		log.Printf("Upload err: %s", err)
		w.WriteHeader(400)
		return
	}

	fmt.Fprint(w, resp.UploadToken)
}
