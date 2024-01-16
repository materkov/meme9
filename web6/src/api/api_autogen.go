package api

import (
	"encoding/json"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/pkg/tracer"
	"github.com/materkov/meme9/web6/src/pkg/xlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strings"
)

var apiCalls = promauto.NewCounter(prometheus.CounterOpts{
	Name: "api_calls",
	Help: "The total number of processed events",
})

func (h *HttpServer) ApiHandler(w http.ResponseWriter, r *http.Request) {
	t := tracer.NewTracer("api")
	defer t.Stop()
	apiCalls.Inc()

	ctx := tracer.WithCtx(r.Context(), t)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Version", pkg.BuildTime)

	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}

	userID := 0
	authHeader := r.Header.Get("authorization")
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader != "" {
		authToken := pkg.ParseAuthToken(ctx, authHeader)
		if authToken != nil {
			userID = authToken.UserID
		}
	}

	viewer := &Viewer{
		UserID:   userID,
		ClientIP: getClientIP(r),
	}

	xlog.Log("Processing API request", xlog.Fields{
		"url":       r.URL.String(),
		"userId":    viewer.UserID,
		"ip":        viewer.ClientIP,
		"userAgent": r.UserAgent(),
	})

	method := strings.TrimPrefix(r.URL.Path, "/api/")
	t.Tags["method"] = method

	switch method {
	case "posts.add":
		req := &PostsAddReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsAdd(ctx, viewer, req)
		writeResp(w, resp, err)

	case "posts.list":
		req := &PostsListReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsList(ctx, viewer, req)
		writeResp(w, resp, err)

	case "posts.delete":
		req := &PostsDeleteReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsDelete(viewer, req)
		writeResp(w, resp, err)

	case "posts.like":
		req := &PostsLikeReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsLike(viewer, req)
		writeResp(w, resp, err)

	case "users.list":
		req := &UsersListReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.usersList(viewer, req)
		writeResp(w, resp, err)

	case "users.setStatus":
		req := &UsersSetStatusReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.usersSetStatus(viewer, req)
		writeResp(w, resp, err)

	case "users.follow":
		req := &UsersFollow{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.usersFollow(viewer, req)
		writeResp(w, resp, err)

	case "auth.login":
		req := &AuthEmailReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.authLogin(viewer, req)
		writeResp(w, resp, err)

	case "auth.register":
		req := &AuthEmailReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.authRegister(viewer, req)
		writeResp(w, resp, err)

	case "auth.vk":
		req := &AuthVkReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.authVk(viewer, req)
		writeResp(w, resp, err)

	case "polls.add":
		req := &PollsAddReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PollsAdd(ctx, viewer, req)
		writeResp(w, resp, err)

	case "polls.list":
		req := &PollsListReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PollsList(ctx, viewer, req)
		writeResp(w, resp, err)

	case "polls.vote":
		req := &PollsVoteReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PollsVote(viewer, req)
		writeResp(w, resp, err)

	case "polls.deleteVote":
		req := &PollsDeleteVoteReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PollsDeleteVote(viewer, req)
		writeResp(w, resp, err)

	default:
		writeResp(w, nil, Error("UnknownMethod"))
	}
}
