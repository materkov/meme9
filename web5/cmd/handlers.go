package main

import (
	"github.com/materkov/meme9/web5/pkg"
	"net/http"
	"strconv"
	"strings"
)

func handleUsersList(w http.ResponseWriter, r *http.Request) {
	ids := parseIds(strings.Split(r.FormValue("ids"), ","))
	viewer := r.Context().Value(ViewerKey).(*Viewer)
	fields := pkg.ParseFields(r.FormValue("fields"))

	postsCursor, _ := strconv.Atoi(r.FormValue("postsCursor"))

	users := usersList(ids, viewer.UserID, fields.Has("isFollowing"), fields.Has("followingCount"), fields.Has("posts"),
		fields.Has("posts.items"),
		fields.Has("posts.items.user"),
		postsCursor,
	)
	write(w, users, nil)
}
