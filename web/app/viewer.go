package app

import (
	"context"

	"github.com/materkov/meme9/web/store"
)

type Viewer struct {
	Token  *store.Token
	UserID int

	RequestHost   string
	RequestScheme string
}

type viewerContextKey struct{}

func GetViewerFromContext(ctx context.Context) *Viewer {
	return ctx.Value(viewerContextKey{}).(*Viewer)
}

func WithViewerContext(parent context.Context, viewer *Viewer) context.Context {
	return context.WithValue(parent, viewerContextKey{}, viewer)
}
