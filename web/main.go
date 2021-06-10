package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/web/pb"
)

func convertPosts(posts []*Post, viewerID int, includeLatestComment bool) []*pb.Post {
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

	photoIdsMap := map[int]bool{}
	for _, post := range posts {
		if post.PhotoID != 0 {
			photoIdsMap[post.PhotoID] = true
		}
	}

	photoIds := make([]int, len(photoIdsMap))
	i = 0
	for photoID := range photoIdsMap {
		photoIds[i] = photoID
		i++
	}

	postIds := make([]int, len(posts))
	for i, post := range posts {
		postIds[i] = post.ID
	}

	likesCountCh := make(chan map[int]int)
	go func() {
		result, err := store.Likes.GetCount(postIds)
		if err != nil {
			log.Printf("Error getting likes count: %s", err)
		}
		likesCountCh <- result
	}()

	isLikedCh := make(chan map[int]bool)
	go func() {
		result, err := store.Likes.GetIsLiked(postIds, viewerID)
		if err != nil {
			log.Printf("Error getting likes count: %s", err)
		}
		isLikedCh <- result
	}()

	usersCh := make(chan []*User)
	go func() {
		result, err := store.User.Get(userIds)
		if err != nil {
			log.Printf("Error selecting users: %s", err)
		}
		usersCh <- result
	}()

	photosCh := make(chan []*Photo)
	go func() {
		result, err := store.Photo.Get(photoIds)
		if err != nil {
			log.Printf("Error selecting photos: %s", err)
		}
		photosCh <- result
	}()

	commentCountsCh := make(chan map[int]int)
	go func() {
		result, err := store.Comment.GetCommentsCounts(postIds)
		if err != nil {
			log.Printf("Error selecting users: %s", err)
		}
		commentCountsCh <- result
	}()

	latestCommentsCh := make(chan map[int]*Comment)
	go func() {
		commentIdsMap, err := store.Comment.GetLatest(postIds)
		if err != nil {
			log.Printf("Error selecting comment ids: %s", err)
		}

		commentIds := make([]int, 0)
		for _, commentID := range commentIdsMap {
			commentIds = append(commentIds, commentID)
		}

		comments, err := store.Comment.Get(commentIds)
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
	photos := <-photosCh

	usersMap := map[int]*User{}
	for _, user := range users {
		usersMap[user.ID] = user
	}

	photosMap := map[int]*Photo{}
	for _, photo := range photos {
		photosMap[photo.ID] = photo
	}

	result := make([]*pb.Post, len(posts))
	for i, post := range posts {
		var wrappedLatestComment *pb.CommentRenderer
		latestComment := latestComments[post.ID]
		if latestComment != nil && includeLatestComment {
			wrappedLatestComment = &pb.CommentRenderer{
				Id:         strconv.Itoa(latestComment.ID),
				Text:       latestComment.Text,
				AuthorId:   strconv.Itoa(latestComment.UserID),
				AuthorName: fmt.Sprintf("User #%d", latestComment.UserID), // TODO
				AuthorUrl:  fmt.Sprintf("/users/%d", latestComment.UserID),
			}
		}

		photoURL := ""
		if post.PhotoID != 0 {
			photo, ok := photosMap[post.PhotoID]
			if ok {
				photoURL = fmt.Sprintf("https://meme-files.s3.eu-central-1.amazonaws.com/photos/%s.jpg", photo.Path)
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
			ImageUrl:      photoURL,
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

var auth *Auth

func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		viewer := Viewer{
			RequestHost: r.Host,
		}

		viewer.RequestScheme = "http"
		if r.Header.Get("x-forwarded-proto") == "https" {
			viewer.RequestScheme = "https"
		}

		viewer.Token, viewer.UserID, _ = auth.tryAuth(r)

		ctx := WithViewerContext(r.Context(), &viewer)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// TODO
var feedSrv *Feed
var postsSrv *Posts
var relationsSrv *Relations
var utilsSrv *Utils
var awsSession *session.Session

func main() {
	rand.Seed(time.Now().UnixNano())
	err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	db, err := sql.Open("mysql", "root:root@/meme9")
	if err != nil {
		log.Fatalf("failed opening mysql connection: %s", err)
	}

	store = NewStore(db)

	auth = &Auth{store: store}

	awsSession, err = session.NewSession(
		&aws.Config{
			Region:      aws.String("eu-central-1"),
			Credentials: credentials.NewStaticCredentials(config.AWSKeyID, config.AWSKeySecret, ""),
		},
	)
	if err != nil {
		log.Fatalf("Failed creating AWS session: %s", err)
	}

	feedSrv = &Feed{}
	utilsSrv = &Utils{}
	postsSrv = &Posts{}
	postsSrv = &Posts{}
	relationsSrv = &Relations{}

	http.Handle("/vk-callback", middleware(handleVKCallback))
	http.Handle("/logout", middleware(handleLogout))
	http.Handle("/upload", middleware(handleUpload))
	http.Handle("/api", middleware(handleAPI))
	http.Handle("/", middleware(handleDefault))

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
