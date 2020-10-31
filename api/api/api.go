package api

import "github.com/materkov/meme9/api/store"

type Viewer struct {
	UserID int
	User   *store.User

	UserAgent string
}
