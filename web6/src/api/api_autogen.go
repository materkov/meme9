package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *HttpServer) ApiHandler(w http.ResponseWriter, r *http.Request) {
	method := strings.TrimPrefix(r.URL.Path, "/api/")
	viewer := r.Context().Value(ctxViewer).(*Viewer)

	switch method {
	case "posts.add":
		req := &PostsAddReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsAdd(viewer, req)
		writeResp(w, resp, err)

	case "posts.list":
		req := &PostsListReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsList(viewer, req)
		writeResp(w, resp, err)

	case "posts.listById":
		req := &PostsListByIdReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsListByID(viewer, req)
		writeResp(w, resp, err)

	case "posts.listByUser":
		req := &PostsListByUserReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PostsListByUser(viewer, req)
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
		req := &UsersSetStatus{}
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
		resp, err := h.Api.PollsAdd(viewer, req)
		writeResp(w, resp, err)

	case "polls.list":
		req := &PollsListReq{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			writeResp(w, nil, ErrParsingRequest)
			return
		}
		resp, err := h.Api.PollsList(viewer, req)
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
