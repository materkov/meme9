package controller

import (
	"github.com/materkov/meme9/web2/store"
	"strconv"
)

type UserPageRenderer struct {
	user  *store.User
	posts []*store.Post
}

func (u *UserPageRenderer) Render() string {
	if u.user == nil {
		return ""
	}

	result := "User page <b>" + u.user.Name + "</b><br/><br/>ID: " + strconv.Itoa(u.user.ID) + "<br/>Name: " + u.user.Name + "<br>"

	result += "<br/><br/>Feed:<br/>"

	for _, post := range u.posts {
		result += PostRenderer{
			post: post,
			user: u.user,
		}.Render()
	}

	return result
}
