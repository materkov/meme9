package api

import "github.com/materkov/meme9/api/store"

type Viewer struct {
	User *store.User

	CSRFValidated bool

	UserAgent     string
	RequestHost   string
	RequestScheme string
}
