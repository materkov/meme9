package types

import (
	"fmt"
	"github.com/materkov/meme9/web4/store"
	"strconv"
)

type User struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Href string `json:"href,omitempty"`
}

func usersList(ids []int) []*User {
	usersMap := map[int]store.User{}
	for _, postID := range ids {
		obj := store.User{}
		err := getObject(postID, &obj)
		if err == nil {
			usersMap[obj.ID] = obj
		}
	}

	results := make([]*User, len(ids))
	for i, userID := range ids {
		result := &User{
			ID:   strconv.Itoa(userID),
			Href: fmt.Sprintf("/users/%d", userID),
		}
		results[i] = result

		user, ok := usersMap[userID]
		if !ok {
			continue
		}

		result.Name = user.Name
	}

	return results
}

type Post struct {
	ID     string `json:"id,omitempty"`
	FromID string `json:"fromId,omitempty"`

	Text       string `json:"text"`
	DetailsURL string `json:"detailsURL"`
}

func postsList(ids []int) []*Post {
	//posts := store.DefaultStore.GetPosts(ids)

	postsMap := map[int]store.Post{}
	for _, postID := range ids {
		post := store.Post{}
		err := getObject(postID, &post)
		if err == nil {
			postsMap[post.ID] = post
		}
	}

	results := make([]*Post, len(ids))
	for i, postID := range ids {
		result := &Post{
			ID:         strconv.Itoa(postID),
			DetailsURL: fmt.Sprintf("/posts/%d", postID),
		}
		results[i] = result

		post, ok := postsMap[postID]
		if !ok {
			continue
		}

		result.FromID = strconv.Itoa(post.UserID)
		result.Text = post.Text
	}

	return results
}
