package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/go-sql-driver/mysql"
	"github.com/materkov/meme9/web/pb"
	"github.com/materkov/meme9/web/store"
)

func convertPosts(ctx context.Context, posts []*store.Post, viewerID int, includeLatestComment bool) []*pb.Post {
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
			count, err := ObjectStore.AssocCount(ctx, postID, store.Assoc_Liked)
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
			data, err := ObjectStore.AssocGet(ctx, postId, store.Assoc_Liked, viewerID)
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

	usersCh := make(chan []*store.User)
	go func() {
		result := make([]*store.User, 0)
		for _, id := range userIds {
			obj, err := ObjectStore.ObjGet(ctx, id)
			if err != nil {
				log.Printf("Error selecting users: %s", err)
				continue
			} else if obj == nil || obj.User == nil {
				log.Printf("User not %d found", id)
				continue
			}

			result = append(result, obj.User)
		}

		usersCh <- result
	}()

	/*photosCh := make(chan []*store.Photo)
	go func() {
		result, err := ObjectStore.ObjGetMany(photoIds)
		if err != nil {
			log.Printf("Error selecting photos: %s", err)
		}

		photos := make([]*store.Photo, 0)
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
			count, err := ObjectStore.AssocCount(ctx, postId, store.Assoc_Commended)
			if err != nil {
				log.Printf("Error selecting users: %s", err)
			} else {
				result[postId] = count
			}
		}

		commentCountsCh <- result
	}()

	latestCommentsCh := make(chan map[int]*store.Comment)
	go func() {
		commentIds := make([]int, 0)
		postLatestComments := map[int]int{}

		for _, postID := range postIds {
			assocs, err := ObjectStore.AssocRange(ctx, postID, store.Assoc_Commended, 1)
			if err != nil {
				log.Printf("Error selecting comment ids: %s", err)
				continue
			}

			if len(assocs) > 0 {
				commentIds = append(commentIds, assocs[0].Commented.ID2)
				postLatestComments[postID] = assocs[0].Commented.ID2
			}
		}

		commentsMap := map[int]*store.Comment{}
		for _, commentID := range commentIds {
			obj, err := ObjectStore.ObjGet(ctx, commentID)
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

		result := map[int]*store.Comment{}
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

	usersMap := map[int]*store.User{}
	for _, user := range users {
		usersMap[user.ID] = user
	}

	//photosMap := map[int]*store.Photo{}
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
			photo, err := ObjectStore.ObjGet(ctx, post.PhotoID)
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

func convertComments(comments []*store.Comment) []*pb.CommentRenderer {
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

// TODO
var FeedSrv *Feed
var postsSrv *Posts
var relationsSrv *Relations
var UtilsSrv *Utils
var awsSession *session.Session

var ObjectStore *store.ObjectStore

func Main() {
	var err error
	awsSession, err = session.NewSession(
		&aws.Config{
			Region:      aws.String("eu-central-1"),
			Credentials: credentials.NewStaticCredentials(DefaultConfig.AWSKeyID, DefaultConfig.AWSKeySecret, ""),
		},
	)
	if err != nil {
		log.Fatalf("Failed creating AWS session: %s", err)
	}

	postsSrv = &Posts{}
	relationsSrv = &Relations{}
}
