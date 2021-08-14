package controller

import (
	"fmt"

	"github.com/materkov/meme9/web2/store"
)

type PostRenderer struct {
	post *store.Post
	user *store.User
}

func (r PostRenderer) Render() string {
	if r.post == nil {
		return ""
	}

	userName := ""
	if r.user != nil {
		userName = r.user.Name
	} else {
		userName = fmt.Sprintf("User #%d", r.post.UserID)
	}

	return fmt.Sprintf("From %s:<br>%s<hr/>", userName, r.post.Text)
}
