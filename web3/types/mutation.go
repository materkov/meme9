package types

import (
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/store"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Mutation struct {
	VKAuth *VKAuth `json:"vkAuth,omitempty"`
}

type VKAuth struct {
	Token string `json:"token,omitempty"`
}

type MutationParams struct {
	AddPost        MutationAddPost        `json:"addPost,omitempty"`
	VKAuthCallback MutationVKAuthCallback `json:"vkAuthCallback,omitempty"`
}

type MutationAddPost struct {
	Include bool   `json:"include,omitempty"`
	Text    string `json:"text,omitempty"`
}

type MutationVKAuthCallback struct {
	Include bool   `json:"include,omitempty"`
	URL     string `json:"url,omitempty"`
}

func ResolveMutation(viewer pkg.Viewer, params MutationParams) *Mutation {
	result := &Mutation{}

	if params.AddPost.Include {
		post := store.Post{
			ID:     store.GenerateID(),
			UserID: viewer.UserID,
			Text:   params.AddPost.Text,
			Date:   int(time.Now().Unix()),
		}

		ch1 := make(chan bool)
		ch2 := make(chan bool)
		go func() {
			_ = GlobalStore.ObjAdd(post.ID, store.ObjectPost, post)
			ch1 <- true
		}()
		go func() {
			_ = GlobalStore.ListAdd(post.UserID, store.ListPosted, post.ID)
			ch2 <- true
		}()
		<-ch1
		<-ch2
	}

	if params.VKAuthCallback.Include {
		urlParsed, err := url.Parse(params.VKAuthCallback.URL)
		if err != nil {
			log.Printf("ERROR")
			return result
		}
		vkID, _ := pkg.ExchangeCode(viewer.Origin, urlParsed.Query().Get("code"))
		log.Printf("UserID %d", vkID)

		userID, _ := GlobalStore.GetMapping(store.MappingVKID, strconv.Itoa(vkID))
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
