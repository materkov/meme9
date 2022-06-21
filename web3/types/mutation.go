package types

import (
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/pkg/globalid"
	"github.com/materkov/web3/store"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Mutation struct {
	Type    string   `json:"type"`
	ID      string   `json:"id"`
	VKAuth  *VKAuth  `json:"vkAuth,omitempty"`
	AddPost *AddPost `json:"addPost,omitempty"`
}

type AddPost struct {
	ID string `json:"id"`
}

type VKAuth struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Token string `json:"token,omitempty"`
}

type MutationParams struct {
	AddPost        *MutationAddPost        `json:"addPost,omitempty"`
	VKAuthCallback *MutationVKAuthCallback `json:"vkAuthCallback,omitempty"`
}

type MutationAddPost struct {
	Text string `json:"text,omitempty"`
}

type MutationVKAuthCallback struct {
	URL string `json:"url,omitempty"`
}

func ResolveMutation(cachedStore *store.CachedStore, viewer pkg.Viewer, params MutationParams) *Mutation {
	result := &Mutation{
		Type: "Mutation",
		ID:   "mutation",
	}

	if params.AddPost != nil {
		post := store.Post{
			ID:     store.GenerateID(),
			UserID: viewer.UserID,
			Text:   params.AddPost.Text,
			Date:   int(time.Now().Unix()),
		}

		ch1 := make(chan bool)
		ch2 := make(chan bool)
		go func() {
			_ = cachedStore.Store.ObjAdd(post.ID, store.ObjectPost, post)
			ch1 <- true
		}()
		go func() {
			_ = cachedStore.Store.ListAdd(post.UserID, store.ListPosted, post.ID)
			ch2 <- true
		}()
		<-ch1
		<-ch2

		result.AddPost = &AddPost{
			ID: globalid.Create(globalid.PostID{PostID: post.ID}),
		}
	}

	if params.VKAuthCallback != nil {
		urlParsed, err := url.Parse(params.VKAuthCallback.URL)
		if err != nil {
			log.Printf("ERROR")
			return result
		}
		vkID, _ := pkg.ExchangeCode(viewer.Origin, urlParsed.Query().Get("code"))
		log.Printf("UserID %d", vkID)

		userID, _ := cachedStore.Store.GetMapping(store.MappingVKID, strconv.Itoa(vkID))
		if userID != 0 {
			token := pkg.AuthToken{
				IssuedAt: int(time.Now().Unix()),
				UserID:   userID,
			}
			result.VKAuth = &VKAuth{Token: token.ToString()}
		}
	}

	return result
}
