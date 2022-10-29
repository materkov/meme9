package contextKeys

type contextKey string

var PostStore = contextKey("PostStore")
var UserStore = contextKey("UserStore")
var LikedStore = contextKey("LikedStore")
