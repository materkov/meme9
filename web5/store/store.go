package store

import (
	"context"
	"github.com/materkov/meme9/web5/pkg/contextKeys"
)

type CachedStore struct {
	Post  GenericCachedStore[Post]
	User  GenericCachedStore[User]
	Photo GenericCachedStore[Photo]

	Liked  LikedStore
	Online OnlineStore
}

func CachedStoreFromCtx(ctx context.Context) *CachedStore {
	return ctx.Value(contextKeys.CachedStore).(*CachedStore)
}

func WithCachedStore(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeys.CachedStore, &CachedStore{
		Post: GenericCachedStore[Post]{
			cache:   map[int]*Post{},
			objType: ObjectTypePost,
		},
		User: GenericCachedStore[User]{
			cache:   map[int]*User{},
			objType: ObjectTypeUser,
		},
		Photo: GenericCachedStore[Photo]{
			cache:   map[int]*Photo{},
			objType: ObjectTypePhoto,
		},
		Liked: LikedStore{
			cache: map[string]likedData{},
		},
		Online: OnlineStore{
			cache:  map[int]bool{},
			needed: map[int]bool{},
		},
	})
}
