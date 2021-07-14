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
	store2 "github.com/materkov/meme9/web/store"
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
		result := map[int]int{}
		for _, postID := range postIds {
			count, err := objectStore.AssocCount(postID, store2.Assoc_Liked)
			if err != nil {
				log.Printf("Error getting likes count: %s", err)
			} else {
				result[postID] = count
			}
		}
		likesCountCh <- result
	}()

	isLikedCh := make(chan map[int]bool)
	go func() {
		result := map[int]bool{}
		for _, postId := range postIds {
			data, err := objectStore.AssocGet(postId, store2.Assoc_Liked, viewerID)
			if err != nil {
				log.Printf("Error getting is Liked: %s", err)
				continue
			}
			if data != nil && data.Liked != nil {
				result[postId] = true
			}
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

	/*photosCh := make(chan []*store2.Photo)
	go func() {
		result, err := objectStore.ObjGetMany(photoIds)
		if err != nil {
			log.Printf("Error selecting photos: %s", err)
		}

		photos := make([]*store2.Photo, 0)
		for _, object := range result {
			if object.Photo != nil {
				photos = append(photos, object.Photo)
			}
		}
		photosCh <- photos
	}()*/

	commentCountsCh := make(chan map[int]int)
	go func() {
		result := map[int]int{}
		for _, postId := range postIds {
			count, err := objectStore.AssocCount(postId, store2.Assoc_Commended)
			if err != nil {
				log.Printf("Error selecting users: %s", err)
			} else {
				result[postId] = count
			}
		}

		commentCountsCh <- result
	}()

	latestCommentsCh := make(chan map[int]*store2.Comment)
	go func() {
		commentIds := make([]int, 0)
		postLatestComments := map[int]int{}

		for _, postID := range postIds {
			assocs, err := objectStore.AssocRange(postID, store2.Assoc_Commended, 1)
			if err != nil {
				log.Printf("Error selecting comment ids: %s", err)
				continue
			}

			if len(assocs) > 0 {
				commentIds = append(commentIds, assocs[0].Commented.ID2)
				postLatestComments[postID] = assocs[0].Commented.ID2
			}
		}

		commentsMap := map[int]*store2.Comment{}
		for _, commentID := range commentIds {
			obj, err := objectStore.ObjGet(commentID)
			if err != nil {
				log.Printf("[Error selecting comment objects: %s", err)
				continue
			}

			if obj == nil || obj.Comment == nil {
				log.Printf("[Error selecting comment objects: empty")
				continue
			}

			commentsMap[obj.Comment.ID] = obj.Comment
		}

		result := map[int]*store2.Comment{}
		for postID, commentID := range postLatestComments {
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
	//photos := <-photosCh

	usersMap := map[int]*User{}
	for _, user := range users {
		usersMap[user.ID] = user
	}

	//photosMap := map[int]*store2.Photo{}
	//for _, photo := range photos {
	//	photosMap[photo.ID] = photo
	//}

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
			photo, err := objectStore.ObjGet(post.PhotoID)
			if err != nil || photo == nil || photo.Photo == nil {
				log.Printf("failed getting photo: %s", err)
			} else {
				photoURL = fmt.Sprintf("https://meme-files.s3.eu-central-1.amazonaws.com/photos/%s.jpg", photo.Photo.Path)
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

func convertComments(comments []*store2.Comment) []*pb.CommentRenderer {
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
	objectStore = store2.NewObjectStore(db)

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
