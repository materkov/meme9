package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/materkov/meme9/web/pb"
)

func convertPosts(posts []*Post, viewerID int) []*pb.Post {
	if len(posts) == 0 {
		return nil
	}

	userIdsMap := map[int]bool{}
	for _, post := range posts {
		userIdsMap[post.UserID] = true
	}

	userIds := make([]int, len(userIdsMap))
	i := 0
	for userID := range userIdsMap {
		userIds[i] = userID
		i++
	}

	postIds := make([]int, len(posts))
	for i, post := range posts {
		postIds[i] = post.ID
	}

	likesCountCh := make(chan map[int]int)
	go func() {
		result, err := store.GetLikesCount(postIds)
		if err != nil {
			log.Printf("Error getting likes count: %s", err)
		}
		likesCountCh <- result
	}()

	isLikedCh := make(chan map[int]bool)
	go func() {
		result, err := store.GetIsLiked(postIds, viewerID)
		if err != nil {
			log.Printf("Error getting likes count: %s", err)
		}
		isLikedCh <- result
	}()

	usersCh := make(chan []*User)
	go func() {
		result, err := store.GetUsers(userIds)
		if err != nil {
			log.Printf("Error selecting users: %s", err)
		}
		usersCh <- result
	}()

	commentCountsCh := make(chan map[int]int)
	go func() {
		result, err := store.GetCommentsCounts(postIds)
		if err != nil {
			log.Printf("Error selecting users: %s", err)
		}
		commentCountsCh <- result
	}()

	latestCommentsCh := make(chan map[int]*Comment)
	go func() {
		commentIdsMap, err := store.GetLatestComments(postIds)
		if err != nil {
			log.Printf("Error selecting comment ids: %s", err)
		}

		commentIds := make([]int, 0)
		for _, commentID := range commentIdsMap {
			commentIds = append(commentIds, commentID)
		}

		comments, err := store.GetComments(commentIds)
		if err != nil {
			log.Printf("[Error selecting comment objects: %s", err)
		}

		commentsMap := map[int]*Comment{}
		for _, comment := range comments {
			commentsMap[comment.ID] = comment
		}

		result := map[int]*Comment{}
		for postID, commentID := range commentIdsMap {
			if comment := commentsMap[commentID]; comment != nil {
				result[postID] = comment
			}
		}

		latestCommentsCh <- result
	}()

	likesCount := <-likesCountCh
	isLiked := <-isLikedCh
	users := <-usersCh
	commentCounts := <-commentCountsCh
	latestComments := <-latestCommentsCh

	usersMap := map[int]*User{}
	for _, user := range users {
		usersMap[user.ID] = user
	}

	result := make([]*pb.Post, len(posts))
	for i, post := range posts {
		var wrappedLatestComment *pb.CommentRenderer
		latestComment := latestComments[post.ID]
		if latestComment != nil {
			wrappedLatestComment = &pb.CommentRenderer{
				Id:         strconv.Itoa(latestComment.ID),
				Text:       latestComment.Text,
				AuthorId:   strconv.Itoa(latestComment.UserID),
				AuthorName: fmt.Sprintf("User #%d", latestComment.UserID), // TODO
				AuthorUrl:  fmt.Sprintf("/users/%d", latestComment.UserID),
			}
		}

		wrappedPost := pb.Post{
			Id:            strconv.Itoa(post.ID),
			Url:           fmt.Sprintf("/posts/%d", post.ID),
			AuthorId:      strconv.Itoa(post.UserID),
			AuthorUrl:     fmt.Sprintf("/users/%d", post.UserID),
			Text:          post.Text,
			DateDisplay:   time.Unix(int64(post.Date), 0).Format("2 Jan 2006 15:04"),
			IsLiked:       isLiked[post.ID],
			LikesCount:    int32(likesCount[post.ID]),
			CanLike:       viewerID != 0,
			CommentsCount: int32(commentCounts[post.ID]),
			TopComment:    wrappedLatestComment,
		}

		user, ok := usersMap[post.UserID]
		if ok {
			wrappedPost.AuthorName = user.Name
			wrappedPost.AuthorAvatar = user.VkAvatar
		}

		result[i] = &wrappedPost
	}

	return result
}

func convertComments(comments []*Comment) []*pb.CommentRenderer {
	result := make([]*pb.CommentRenderer, len(comments))
	for i, comment := range comments {
		result[i] = &pb.CommentRenderer{
			Id:         strconv.Itoa(comment.ID),
			Text:       comment.Text,
			AuthorId:   strconv.Itoa(comment.UserID),
			AuthorName: fmt.Sprintf("user %d", comment.UserID),
			AuthorUrl:  fmt.Sprintf("/users/%d", comment.UserID),
		}
	}

	return result
}

func fetchVkData(userId int, accessToken string) (string, string, error) {
	resp, err := http.PostForm("https://api.vk.com/method/users.get", url.Values{
		"access_token": []string{accessToken},
		"v":            []string{"5.130"},
		"user_ids":     []string{strconv.Itoa(userId)},
		"fields":       []string{"photo_200,first_name,last_name"},
	})
	if err != nil {
		return "", "", fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading http body: %w", err)
	}

	body := struct {
		Response []struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Photo200  string `json:"photo_200"`
		}
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return "", "", fmt.Errorf("failed parsing json: %w", err)
	}

	if len(body.Response) == 0 {
		return "", "", fmt.Errorf("response length zero: %s", bodyBytes)
	}

	return body.Response[0].FirstName + " " + body.Response[0].LastName, body.Response[0].Photo200, nil
}

func twirpWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		viewer := Viewer{
			RequestHost:   r.Host,
			RequestScheme: "http",
		}

		accessCookie, err := r.Cookie("access_token")
		if err == nil && accessCookie.Value != "" {
			token, err := store.GetToken(accessCookie.Value)
			if err == nil {
				viewer.Token = token
				viewer.UserID = token.UserID
			}
		}

		if r.Header.Get("x-forwarded-proto") == "https" {
			viewer.RequestScheme = "https"
		}

		ctx := WithViewerContext(r.Context(), &viewer)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TODO
var feedSrv *Feed
var utilsSrv *Utils

func main() {
	rand.Seed(time.Now().UnixNano())
	config.MustLoad()

	db, err := sqlx.Open("mysql", "root:root@/meme9")
	if err != nil {
		panic(err)
	}

	store = Store{db: db}

	feedSrv = &Feed{}
	utilsSrv = &Utils{}

	// Twirp API
	http.Handle("/twirp/meme.Feed/", twirpWrapper(pb.NewFeedServer(feedSrv)))
	http.Handle("/twirp/meme.Profile/", twirpWrapper(pb.NewProfileServer(&Profile{})))
	http.Handle("/twirp/meme.Relations/", twirpWrapper(pb.NewRelationsServer(&Relations{})))
	http.Handle("/twirp/meme.Posts/", twirpWrapper(pb.NewPostsServer(&Posts{})))
	http.Handle("/twirp/meme.Utils/", twirpWrapper(pb.NewUtilsServer(utilsSrv)))

	// Other
	http.Handle("/vk-callback", twirpWrapper(http.HandlerFunc(handleVKCallback)))
	http.Handle("/logout", twirpWrapper(http.HandlerFunc(handleLogout)))
	http.Handle("/", twirpWrapper(http.HandlerFunc(handleDefault)))

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
