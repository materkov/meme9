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

type PostPageResponse struct {
	Route    Route  `json:"route,omitempty"`
	PagePost string `json:"pagePost,omitempty"`
	Nodes    *Nodes `json:"nodes,omitempty"`
}

func PostPage(req *PostPageRequest) (*PostPageResponse, error) {
	postID, _ := strconv.Atoi(req.PostID)

	posts := postsList([]int{postID})
	users := usersList(getUsersFromPosts(posts))

	return &PostPageResponse{
		Route:    RoutePostsId,
		PagePost: req.PostID,
		Nodes: &Nodes{
			Posts: posts,
			Users: users,
		},
	}, nil
}

type UserPageRequest struct {
	UserID string `json:"userId,omitempty"`
}

type UserPageResponse struct {
	Route    Route    `json:"route,omitempty"`
	PageUser string   `json:"pageUser,omitempty"`
	NotFound bool     `json:"notFound,omitempty"`
	Nodes    *Nodes   `json:"nodes,omitempty"`
	Posts    []string `json:"posts,omitempty"`
}

func UserPage(req *UserPageRequest) (*UserPageResponse, error) {
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

	return &UserPageResponse{
		Route:    RouteUserPage,
		PageUser: strconv.Itoa(user.ID),
		Posts:    postIdsStr,
		Nodes: &Nodes{
			Users: usersList([]int{userID}),
			Posts: postsList(postIds),
		},
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

type FeedResponse struct {
	Route Route    `json:"route,omitempty"`
	Posts []string `json:"posts,omitempty"`
	Nodes *Nodes   `json:"nodes,omitempty"`
}

func Feed(req *FeedRequest) (*FeedResponse, error) {
	postIds, err := postsGetFeed()
	if err != nil {
		log.Printf("Error getting feed: %s", err)
	}

	posts := postsList(postIds)
	users := usersList(getUsersFromPosts(posts))

	postIdsStr := make([]string, len(postIds))
	for i, postID := range postIds {
		postIdsStr[i] = strconv.Itoa(postID)
	}

	return &FeedResponse{
		Posts: postIdsStr,
		Nodes: &Nodes{
			Posts: posts,
			Users: users,
		},
	}, nil
}
