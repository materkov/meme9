package handlers

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers/common"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

type GetFeed struct {
	Store *store.Store
}

func (g *GetFeed) Handle(viewer *api.Viewer, req *pb.GetFeedRequest) (*pb.GetFeedRenderer, error) {
	postIds, err := g.Store.GetFeed()
	if err != nil {
		log.Printf("[ERROR] Error getting feed: %s", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(postIds))

	posts := make([]*store.Post, len(postIds))
	for i := range postIds {
		idxCopy := i
		go func() {
			post, err := g.Store.GetPost(postIds[idxCopy])
			if err != nil {
				log.Printf("[ERROR] Error getting post for feed: %s", err)
			}

			posts[idxCopy] = post
			wg.Done()
		}()
	}
	wg.Wait()

	postPageRenderers := make([]*pb.PostPageRenderer, 0)
	for i, post := range posts {
		if posts[i] == nil {
			continue
		}

		postPageRenderers = append(postPageRenderers, &pb.PostPageRenderer{
			Id:      strconv.Itoa(post.ID),
			PostUrl: fmt.Sprintf("/posts/%d", post.ID),
			Text:    post.Text,
			UserId:  strconv.Itoa(post.UserID),
			UserUrl: fmt.Sprintf("/users/%d", post.UserID),
		})
	}

	renderer := &pb.GetFeedRenderer{
		Posts:          postPageRenderers,
		HeaderRenderer: common.GetHeaderRenderer(viewer),
	}

	return renderer, nil
}
