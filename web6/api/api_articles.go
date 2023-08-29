package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/materkov/meme9/web6/pkg"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//go:generate easyjson -all -lower_camel_case -omit_empty api_articles.go

type apiHandler func(w http.ResponseWriter, r *http.Request, token *pkg.AuthToken) (interface{}, error)

func wrapAPI(handler apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Version", pkg.BuildTime)

		var authToken *pkg.AuthToken

		authHeader := r.Header.Get("authorization")
		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader != "" {
			authToken = pkg.ParseAuthToken(authHeader)
		}

		resp, err := handler(w, r, authToken)
		if err != nil {
			w.WriteHeader(400)
			var publicErr *Error
			if ok := errors.As(err, &publicErr); ok {
				fmt.Fprint(w, publicErr.Message)
			} else {
				fmt.Fprint(w, "Internal server error")
			}
			return
		}

		_ = json.NewEncoder(w).Encode(resp)
	}
}

type webHandler func(w http.ResponseWriter, r *http.Request, token *pkg.AuthToken)

func wrapWeb(handler webHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Version", pkg.BuildTime)

		var authToken *pkg.AuthToken

		authCookie, _ := r.Cookie("authToken")
		if authCookie != nil {
			authToken = pkg.ParseAuthToken(authCookie.Value)
		}

		handler(w, r, authToken)
	}
}

type Void struct{}

func (*HttpServer) usersList(w http.ResponseWriter, r *http.Request, t *pkg.AuthToken) (interface{}, error) {
	req := struct {
		UserIds []string
	}{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, &Error{400, "error parsing request"}
	}

	result := make([]*User, len(req.UserIds))
	for i, userIdStr := range req.UserIds {
		userId, _ := strconv.Atoi(userIdStr)
		user, err := pkg.GetUser(userId)
		if err != nil {
			log.Printf("[ERROR] Error loading user: %s", err)
		}

		result[i] = transformUser(userId, user)
	}

	return result, nil
}
