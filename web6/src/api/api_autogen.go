package api

import (
	"encoding/json"
	"net/http"
)

func (h *HttpServer) PostsAdd(w http.ResponseWriter, r *http.Request) {
	req := &PostsAddReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.PostsAdd(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) PostsList(w http.ResponseWriter, r *http.Request) {
	req := &PostsListReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.PostsList(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) PostsListByID(w http.ResponseWriter, r *http.Request) {
	req := &PostsListByIdReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.PostsListByID(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) PostsListByUser(w http.ResponseWriter, r *http.Request) {
	req := &PostsListByUserReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.PostsListByUser(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) PostsDelete(w http.ResponseWriter, r *http.Request) {
	req := &PostsDeleteReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.PostsDelete(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) usersList(w http.ResponseWriter, r *http.Request) {
	req := &UsersListReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.usersList(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) usersSetStatus(w http.ResponseWriter, r *http.Request) {
	req := &UsersSetStatus{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.usersSetStatus(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) usersFollow(w http.ResponseWriter, r *http.Request) {
	req := &UsersFollow{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.usersFollow(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) authLogin(w http.ResponseWriter, r *http.Request) {
	req := &AuthEmailReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.authLogin(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) authRegister(w http.ResponseWriter, r *http.Request) {
	req := &AuthEmailReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.authRegister(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}

func (h *HttpServer) authVK(w http.ResponseWriter, r *http.Request) {
	req := &AuthVkReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeResp(w, nil, ErrParsingRequest)
		return
	}
	resp, err := h.Api.authVk(r.Context().Value(ctxViewer).(*Viewer), req)
	writeResp(w, resp, err)
}
