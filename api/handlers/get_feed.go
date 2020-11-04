package handlers

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers/common"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type GetFeed struct {
	Store *store.Store
}

func (g *GetFeed) Handle(viewer *api.Viewer, req *login.GetFeedRequest) *login.AnyRenderer {
	postIds, err := g.Store.GetFeed()
	if err != nil {
		return &login.AnyRenderer{Renderer: &login.AnyRenderer_ErrorRenderer{
			ErrorRenderer: &login.ErrorRenderer{
				DisplayText: "eer",
			},
		}}
	}

	wg := sync.WaitGroup{}
	wg.Add(len(postIds))

	posts := make([]*store.Post, len(postIds))
	for i := range postIds {
		idxCopy := i
		go func() {
			post, _ := g.Store.GetPost(postIds[idxCopy])
			posts[idxCopy] = post
			wg.Done()
		}()
	}
	wg.Wait()

	postPageRenderers := make([]*login.PostPageRenderer, len(posts))
	for i, post := range posts {
		postPageRenderers[i] = &login.PostPageRenderer{
			Id:            strconv.Itoa(post.ID),
			PostUrl:       fmt.Sprintf("/posts/%d", post.ID),
			Text:          post.Text,
			UserId:        strconv.Itoa(post.UserID),
			UserUrl:       fmt.Sprintf("/users/%d", post.UserID),
		}
	}

	return &login.AnyRenderer{Renderer: &login.AnyRenderer_GetFeedRenderer{
		GetFeedRenderer: &login.GetFeedRenderer{
			Posts:          postPageRenderers,
			HeaderRenderer: common.GetHeaderRenderer(viewer),
		},
	}}
}
