package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/materkov/meme9/web/pb"
)

func handleIndex(_ string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	if viewer.UserID == 0 {
		return &pb.UniversalRenderer{
			Renderer: &pb.UniversalRenderer_FeedRenderer{FeedRenderer: &pb.FeedRenderer{}},
		}, nil
	}

	following, err := store.GetFollowing(viewer.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting following ids: %w", err)
	}

	following = append(following, viewer.UserID)

	postIds, err := store.GetPostsByUsers(following)
	if err != nil {
		return nil, fmt.Errorf("error getting post ids: %w", err)
	}

	posts, err := store.GetPosts(postIds)
	if err != nil {
		return nil, fmt.Errorf("error getting post ids: %w", err)
	}

	wrappedPosts := convertPosts(posts, viewer.UserID)

	return &pb.UniversalRenderer{
		Renderer: &pb.UniversalRenderer_FeedRenderer{FeedRenderer: &pb.FeedRenderer{
			Posts: wrappedPosts,
		}},
	}, nil
}

func handleLogin(_ string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	requestScheme := viewer.RequestScheme
	requestHost := viewer.RequestHost
	vkAppID := 7260220
	redirectURL := url.QueryEscape(fmt.Sprintf("%s://%s/vk-callback", requestScheme, requestHost))
	vkURL := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", vkAppID, redirectURL)

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_LoginPageRenderer{
		LoginPageRenderer: &pb.LoginPageRenderer{
			AuthUrl: vkURL,
			Text:    "Войти через ВК",
		},
	}}, nil
}

func handleProfile(url string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	req := &pb.ProfileGetRequest{
		Id: strings.TrimPrefix(url, "/users/"),
	}

	userID, _ := strconv.Atoi(req.Id)
	users, err := store.GetUsers([]int{userID})
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %w", err)
	} else if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := users[0]

	postIds, err := store.GetPostsByUsers([]int{userID})
	if err != nil {
		log.Printf("Error selecting user posts: %s", err)
	}

	posts, err := store.GetPosts(postIds)
	if err != nil {
		log.Printf("Error selecting posts: %s", err)
	}

	wrappedPosts := convertPosts(posts, viewer.UserID)

	followingIds, err := store.GetFollowing(viewer.UserID)
	if err != nil {
		log.Printf("Error getting following users: %s", err)
	}

	isFollowing := false
	for _, userID := range followingIds {
		if userID == user.ID {
			isFollowing = true
		}
	}

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_ProfileRenderer{ProfileRenderer: &pb.ProfileRenderer{
		Id:          strconv.Itoa(user.ID),
		Name:        user.Name,
		Avatar:      user.VkAvatar,
		Posts:       wrappedPosts,
		IsFollowing: isFollowing,
	}}}, nil
}

func handlePostPage(url string, viewer *Viewer) (*pb.UniversalRenderer, error) {
	postIDStr := strings.TrimPrefix(url, "/posts/")
	postID, _ := strconv.Atoi(postIDStr)

	posts, err := store.GetPosts([]int{postID})
	if err != nil {
		return nil, fmt.Errorf("error selecting post: %s", err)
	} else if len(posts) == 0 {
		return nil, fmt.Errorf("post not found")
	}

	commentIds, err := store.GetCommentsByPost(postID)
	if err != nil {
		log.Printf("Error selecting comment ids: %s", err)
	}

	comments, err := store.GetComments(commentIds)
	if err != nil {
		log.Printf("Error selecting comments objects: %s", err)
	}

	// TODO
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].ID > comments[j].ID
	})

	wrappedPosts := convertPosts(posts, viewer.UserID)
	wrappedComments := convertComments(comments)

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_PostRenderer{PostRenderer: &pb.PostRenderer{
		Post:     wrappedPosts[0],
		Comments: wrappedComments,
		Composer: &pb.CommentComposerRenderer{
			PostId:      postIDStr,
			Placeholder: "Напишите здесь свой комментарий...",
		},
	}}}, nil
}

func convertPosts(posts []*Post, viewerID int) []*pb.Post {
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

func handleVKCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		_, _ = fmt.Fprint(w, "Empty VK code")
		return
	}

	// TODO
	proxyScheme := "http"
	if r.Header.Get("x-forwarded-proto") != "" {
		proxyScheme = r.Header.Get("x-forwarded-proto")
	}

	vkAppID := 7260220

	redirectURI := fmt.Sprintf("%s://%s%s", proxyScheme, r.Host, r.URL.EscapedPath())

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(vkAppID)},
		"client_secret": []string{config.VKAppSecret},
		"redirect_uri":  []string{redirectURI},
		"code":          []string{code},
	})
	if err != nil {
		fmt.Fprintf(w, "http vk error: %v", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "failed reading http body: %v", err)
		return
	}

	body := struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		fmt.Fprintf(w, "incorrect json: %s", bodyBytes)
		return
	} else if body.AccessToken == "" {
		fmt.Fprintf(w, "incorrect response: %s", bodyBytes)
		return
	}

	userID, err := store.GetByVkID(body.UserID)
	if err != nil {
		log.Printf("Error selecting by vk id: %s", err)
		fmt.Fprintf(w, "internal error")
		return
	}

	users, err := store.GetUsers([]int{userID})
	if err != nil {
		log.Printf("Error selecting user: %s", err)
		fmt.Fprintf(w, "internal error")
		return
	} else if len(users) == 0 {
		log.Printf("User %d not found", userID)
		fmt.Fprintf(w, "internal error")
		return
	}

	user := users[0]

	vkName, vkAvatar, err := fetchVkData(body.UserID, body.AccessToken)
	if err == nil {
		user.Name = vkName
		user.VkAvatar = vkAvatar
		err = store.UpdateNameAvatar(user)
		if err != nil {
			log.Printf("Failed saving new name&avatar: %s", err)
		}
	}

	token := Token{
		Token:  RandString(50),
		UserID: userID,
	}
	err = store.AddToken(&token)
	if err != nil {
		log.Printf("error saving token: %s", err)
		fmt.Fprintf(w, "internal error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.Token,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
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

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func twirpWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")

		viewer := Viewer{}
		accessCookie, err := r.Cookie("access_token")
		if err == nil && accessCookie.Value != "" {
			token, err := store.GetToken(accessCookie.Value)
			if err == nil {
				viewer.Token = token
				viewer.UserID = token.UserID
			}
		}

		viewer.RequestScheme = "http"
		if r.Header.Get("x-forwarded-proto") != "" {
			viewer.RequestScheme = r.Header.Get("x-forwarded-proto")
		}

		viewer.RequestHost = r.Host

		ctx := WithViewerContext(r.Context(), &viewer)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	configStr, _ := os.LookupEnv("CONFIG")
	_ = json.Unmarshal([]byte(configStr), &config)

	db, err := sqlx.Open("mysql", "root:root@/meme9")
	if err != nil {
		panic(err)
	}

	store = Store{db: db}

	// Twirp API
	http.Handle("/twirp/meme.Feed/", twirpWrapper(pb.NewFeedServer(&Feed{})))
	http.Handle("/twirp/meme.Profile/", twirpWrapper(pb.NewProfileServer(&Profile{})))
	http.Handle("/twirp/meme.Relations/", twirpWrapper(pb.NewRelationsServer(&Relations{})))
	http.Handle("/twirp/meme.Posts/", twirpWrapper(pb.NewPostsServer(&Posts{})))
	http.Handle("/twirp/meme.Utils/", twirpWrapper(pb.NewUtilsServer(&Utils{})))

	// Other
	http.HandleFunc("/vk-callback", handleVKCallback)
	http.HandleFunc("/logout", handleLogout)

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
