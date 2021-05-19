package main

import "context"

type Viewer struct {
	Token  *Token
	UserID int

	RequestHost   string
	RequestScheme string
}

func GetViewerFromContext(ctx context.Context) *Viewer {
	return ctx.Value("viewer").(*Viewer)
}

func WithViewerContext(parent context.Context, viewer *Viewer) context.Context {
	return context.WithValue(parent, "viewer", viewer)
}
