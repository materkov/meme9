package types

import (
	"fmt"
	"github.com/materkov/meme9/web4/store"
	"log"
	"strconv"
)

type PostPageRequest struct {
	PostID string `json:"postId,omitempty"`
}

func PostPage(req *PostPageRequest) ([]interface{}, error) {
	postID, _ := strconv.Atoi(req.PostID)

	posts := postsList([]int{postID})
	users := getUsersFromPosts(posts)
	usersList(users)

	return []interface{}{
		req.PostID,
		posts[0],
	}, nil
}

type UserPageRequest struct {
	UserID string `json:"userId,omitempty"`
}

func UserPage(req *UserPageRequest) ([]interface{}, error) {
	userID, _ := strconv.Atoi(req.UserID)
	if userID <= 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := &store.User{}
	err := getObject(userID, &user)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	postIds, err := postsGetFeedByUsers([]int{user.ID})
	if err != nil {
		log.Printf("Error getting user feed: %s", err)
	}

	postIdsStr := make([]string, len(postIds))
	for i, postID := range postIds {
		postIdsStr[i] = strconv.Itoa(postID)
	}

	users := []*User{{ID: strconv.Itoa(userID)}}
	usersList(users)

	return []interface{}{
		users[0],
		postsList(postIds),
	}, nil
}

type VkCallbackRequest struct {
	Code string
}

type VkCallbackResponse struct {
	UserID    string `json:"userId,omitempty"`
	AuthToken string `json:"authToken,omitempty"`
}

func VkCallback(req *VkCallbackRequest) (*VkCallbackResponse, error) {
	vkID, _ := authExchangeCode("http://localhost:3000", req.Code)
	if vkID == 0 {
		return nil, fmt.Errorf("error exchanging code")
	}

	userID, err := usersGetOrCreateByVKID(vkID)
	if err != nil {
		return nil, fmt.Errorf("error exchanging code")
	}

	token, err := authCreateToken(userID)
	if err != nil {
		return nil, fmt.Errorf("error creating token")
	}

	return &VkCallbackResponse{
		UserID:    strconv.Itoa(userID),
		AuthToken: token,
	}, nil
}

type AddPostRequest struct {
	Text string `json:"text"`
}

type AddPostResponse struct {
	Post *Post `json:"post,omitempty"`
}

func AddPost(req *AddPostRequest, viewer *Viewer) (*AddPostResponse, error) {
	postID, _ := postsAdd(req.Text, viewer)
	post := postsList([]int{postID})[0]
	return &AddPostResponse{Post: post}, nil
}

type FeedRequest struct {
}

func Feed(req *FeedRequest) ([]interface{}, error) {
	postIds, err := postsGetFeed()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	posts := postsList(postIds)
	users := getUsersFromPosts(posts)
	usersList(users)

	postIdsStr := make([]string, len(postIds))
	for i, postID := range postIds {
		postIdsStr[i] = strconv.Itoa(postID)
	}

	return []interface{}{
		posts,
	}, nil
}
