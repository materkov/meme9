package types

import (
	"github.com/materkov/web3/store"
	"time"
)

type Mutation struct {
}

type MutationParams struct {
	AddPost MutationAddPost `json:"addPost,omitempty"`
}

type MutationAddPost struct {
	Include bool   `json:"include"`
	Text    string `json:"text"`
}

func ResolveMutation(params MutationParams) *Mutation {
	result := &Mutation{}

	if params.AddPost.Include {
		post := store.Post{
			ID:     store.GenerateID(),
			UserID: 10,
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

	return result
}
