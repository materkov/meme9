package contextKeys

type contextKey string

var PostStore = contextKey("PostStore")
var UserStore = contextKey("UserStore")
var LikedStore = contextKey("LikedStore")
var OnlineStore = contextKey("OnlineStore")
var CachedStore = contextKey("CachedStore")
