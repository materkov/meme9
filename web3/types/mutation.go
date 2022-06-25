package types

import (
	"fmt"
	"github.com/materkov/web3/pkg"
	"github.com/materkov/web3/pkg/globalid"
	"github.com/materkov/web3/store"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Mutation struct {
	Type     string    `json:"type"`
	ID       string    `json:"id"`
	VKAuth   *VKAuth   `json:"vkAuth,omitempty"`
	AddPost  *AddPost  `json:"addPost,omitempty"`
	Follow   *Follow   `json:"follow"`
	Unfollow *Unfollow `json:"unfollow"`
}

type AddPost struct {
	ID string `json:"id"`
}

type MutationFollow struct {
	UserID string `json:"userId"`
}

type Follow struct {
}

type MutationUnfollow struct {
	UserID string `json:"userId"`
}

type Unfollow struct {
}

type VKAuth struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Token string `json:"token,omitempty"`
}

type MutationParams struct {
	AddPost        *MutationAddPost        `json:"addPost,omitempty"`
	VKAuthCallback *MutationVKAuthCallback `json:"vkAuthCallback,omitempty"`
	Follow         *MutationFollow         `json:"follow,omitempty"`
	Unfollow       *MutationUnfollow       `json:"unfollow,omitempty"`
}

type MutationAddPost struct {
	Text string `json:"text,omitempty"`
}

type MutationVKAuthCallback struct {
	URL string `json:"url,omitempty"`
}

func ResolveMutation(cachedStore *store.CachedStore, viewer *pkg.Viewer, params MutationParams) *Mutation {
	result := &Mutation{
		Type: "Mutation",
		ID:   "mutation",
	}

	if params.AddPost != nil {
		post := &store.Post{
			ID:     store.GenerateID(),
			UserID: viewer.UserID,
			Text:   params.AddPost.Text,
			Date:   int(time.Now().Unix()),
		}

		ch1 := make(chan bool)
		ch2 := make(chan bool)
		go func() {
			_ = cachedStore.Store.ObjAdd(store.ObjectPost, post)
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
		if userID == 0 {
			user := &store.User{
				ID:   store.GenerateID(),
				Name: fmt.Sprintf("User #%d", vkID),
				VkID: vkID,
			}
			_ = cachedStore.Store.ObjAdd(store.ObjectUser, user)

			userID = user.ID

			err = cachedStore.Store.SaveMapping(store.MappingVKID, strconv.Itoa(vkID), user.ID)
			if err != nil {
				log.Printf("erER")
				return result
			}
		}

		token := pkg.AuthToken{
			IssuedAt: int(time.Now().Unix()),
			UserID:   userID,
		}
		result.VKAuth = &VKAuth{Token: token.ToString()}
	}

	if params.Follow != nil {
		if viewer.UserID != 0 {
			userID, _ := globalid.Parse(params.Follow.UserID)
			if userID, ok := userID.(*globalid.UserID); ok {
				_, err := pkg.GlobalStore.ListGetItem(viewer.UserID, store.ListSubscribedTo, userID.UserID)

				if err == store.ErrListItemNotExists {
					_ = pkg.GlobalStore.ListAdd(viewer.UserID, store.ListSubscribedTo, userID.UserID)
				}
			}
		}
		result.Follow = &Follow{}
	}

	if params.Unfollow != nil {
		if viewer.UserID != 0 {
			userID, _ := globalid.Parse(params.Unfollow.UserID)
			if userID, ok := userID.(*globalid.UserID); ok {
				_ = pkg.GlobalStore.ListDel(viewer.UserID, store.ListSubscribedTo, userID.UserID)
			}
		}
	}

	return result
}
