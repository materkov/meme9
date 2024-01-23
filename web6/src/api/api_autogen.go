package api

import (
	"github.com/materkov/meme9/web6/src/pkg"
	"io"
	"net/http"
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

	currentMethod := r.URL.Path

	var methodsMap = map[string]string{
		"/api/posts.add":        "/twirp/meme.api.Posts/Add",
		"/api/posts.list":       "/twirp/meme.api.Posts/List",
		"/api/posts.delete":     "/twirp/meme.api.Posts/Delete",
		"/api/posts.like":       "/twirp/meme.api.Posts/Like",
		"/api/users.list":       "/twirp/meme.api.Users/List",
		"/api/users.setStatus":  "/twirp/meme.api.Users/SetStatus",
		"/api/users.follow":     "/twirp/meme.api.Users/Follow",
		"/api/auth.login":       "/twirp/meme.api.Auth/Login",
		"/api/auth.register":    "/twirp/meme.api.Auth/Register",
		"/api/auth.vk":          "/twirp/meme.api.Auth/Vk",
		"/api/polls.add":        "/twirp/meme.api.Polls/Add",
		"/api/polls.list":       "/twirp/meme.api.Polls/List",
		"/api/polls.vote":       "/twirp/meme.api.Polls/Vote",
		"/api/polls.deleteVote": "/twirp/meme.api.Polls/DeleteVote",
	}

	apiReq, _ := http.NewRequest("POST", "http://localhost:8002"+methodsMap[currentMethod], r.Body)
	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Authorization", r.Header.Get("Authorization"))

	resp, err := http.DefaultClient.Do(apiReq)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	_, _ = io.Copy(w, resp.Body)
	_ = resp.Body.Close()
}
