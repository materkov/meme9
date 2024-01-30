package server

import "github.com/twitchtv/twirp"

var ErrNotAuthorized = twirp.NewError(twirp.InvalidArgument, "NotAuthorized")
