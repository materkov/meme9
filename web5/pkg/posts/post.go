package posts

import "github.com/materkov/meme9/web5/store"

func CanSee(post *store.Post, viewerID int) bool {
	if post == nil {
		return false
	}
	if post.IsDeleted {
		return false
	}

	return true
}

func CanEdit(post *store.Post, viewerID int) bool {
	if !CanSee(post, viewerID) {
		return false
	}

	return post.UserID == viewerID
}
