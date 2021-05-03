package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/materkov/meme9/web/pb"
)

type apiHandler func(body io.Reader, viewer *Viewer) (proto.Message, error)

func apiHandlerFeedGetHeader(body io.Reader, viewer *Viewer) (proto.Message, error) {
	headerRenderer := pb.HeaderRenderer{
		MainUrl:   "/",
		LogoutUrl: "http://localhost:8000/logout",
	}

	if viewer.UserID != 0 {
		users, err := store.GetUsers([]int{viewer.UserID})
		if err != nil {
			log.Printf("Error getting user: %s", err)
		} else if len(users) == 0 {
			log.Printf("User %d not found", viewer.UserID)
		} else {
			user := users[0]

			headerRenderer.IsAuthorized = true
			headerRenderer.UserAvatar = user.VkAvatar
			headerRenderer.UserName = user.Name
		}
	}

	return &pb.FeedGetHeaderResponse{Renderer: &headerRenderer}, nil
}

func apiHandlerPostsAdd(body io.Reader, viewer *Viewer) (proto.Message, error) {
	req := pb.PostsAddRequest{}
	err := jsonpb.Unmarshal(body, &req)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling request: %w", err)
	}

	post := Post{
		UserID: viewer.UserID,
		Date:   int(time.Now().Unix()),
		Text:   req.Text,
	}

	err = store.AddPost(&post)
	if err != nil {
		return nil, fmt.Errorf("error saving post: %w", err)
	}

	return &pb.PostsAddResponse{
		PostUrl: fmt.Sprintf("/profile/%d", viewer.UserID),
	}, nil
}

func apiHandlerResolveRoute(body io.Reader, viewer *Viewer) (proto.Message, error) {
	request := pb.ResolveRouteRequest{}
	err := jsonpb.Unmarshal(body, &request)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling request: %w", err)
	}

	if request.Url == "/" {
		return handleIndex(request.Url)
	} else if request.Url == "/login" {
		return handleIndex(request.Url)
	} else if m, _ := regexp.MatchString(`^/users/\d+$`, request.Url); m {
		return handleProfile(request.Url)
	} else if m, _ := regexp.MatchString(`^/posts/\d+$`, request.Url); m {
		return handleIndex(request.Url)
	} else {
		return &pb.UniversalRenderer{}, nil
	}
}

type routeHandler func(url string) (*pb.UniversalRenderer, error)

func handleIndex(_ string) (*pb.UniversalRenderer, error) {
	postIds, err := store.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("error getting post ids: %w", err)
	}

	posts, err := store.GetPosts(postIds)
	if err != nil {
		return nil, fmt.Errorf("error getting post ids: %w", err)
	}

	wrappedPosts := convertPosts(posts)

	return &pb.UniversalRenderer{
		Renderer: &pb.UniversalRenderer_FeedRenderer{FeedRenderer: &pb.FeedRenderer{
			Posts: wrappedPosts,
		}},
	}, nil
}

func handleLogin(_ string) (*pb.UniversalRenderer, error) {
	requestScheme := "http"
	requestHost := "localhost:8000"
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

func handleProfile(url string) (*pb.UniversalRenderer, error) {
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

	postIds, err := store.GetPostsByUser(userID)
	if err != nil {
		log.Printf("Error selecting user posts: %s", err)
	}

	posts, err := store.GetPosts(postIds)
	if err != nil {
		log.Printf("Error selecting posts: %s", err)
	}

	wrappedPosts := convertPosts(posts)

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_ProfileRenderer{ProfileRenderer: &pb.ProfileRenderer{
		Id:     strconv.Itoa(user.ID),
		Name:   user.Name,
		Avatar: user.VkAvatar,
		Posts:  wrappedPosts,
	}}}, nil
}

func handlePostPage(url string) (*pb.UniversalRenderer, error) {
	postIDStr := strings.TrimPrefix(url, "/posts/")
	postID, _ := strconv.Atoi(postIDStr)

	posts, err := store.GetPosts([]int{postID})
	if err != nil {
		return nil, fmt.Errorf("error selecting post: %s", err)
	} else if len(posts) == 0 {
		return nil, fmt.Errorf("post not found")
	}

	wrappedPosts := convertPosts(posts)

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_PostRenderer{PostRenderer: &pb.PostRenderer{
		Post: wrappedPosts[0],
	}}}, nil
}

func convertPosts(posts []*Post) []*pb.Post {
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

	users, err := store.GetUsers(userIds)
	if err != nil {
		log.Printf("Error selecting users: %s", err)
	}

	usersMap := map[int]*User{}
	for _, user := range users {
		usersMap[user.ID] = user
	}

	result := make([]*pb.Post, len(posts))
	for i, post := range posts {

		wrappedPost := pb.Post{
			Id:          strconv.Itoa(post.ID),
			Url:         fmt.Sprintf("/posts/%d", post.ID),
			AuthorId:    strconv.Itoa(post.UserID),
			AuthorUrl:   fmt.Sprintf("/users/%d", post.UserID),
			Text:        post.Text,
			DateDisplay: time.Unix(int64(post.Date), 0).Format("2 Jan 2006 15:04"),
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

func handleVKCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		_, _ = fmt.Fprint(w, "Empty VK code")
		return
	}

	proxyScheme := "http"
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
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.Token,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "http://localhost:3000", http.StatusFound)
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

	http.Redirect(w, r, "http://localhost:3000", http.StatusFound)
}

func apiWrapper(next apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		resp, err := next(r.Body, &viewer)
		if err != nil {
			log.Printf("Error response: %s", err)

			w.WriteHeader(400)
			_ = json.NewEncoder(w).Encode(struct {
				Error string `json:"error"`
			}{
				Error: "oops, error",
			})
			return
		}

		m := jsonpb.Marshaler{}
		_ = m.Marshal(w, resp)
	}
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

	r := mux.NewRouter()

	// API
	r.HandleFunc("/api/meme.Posts.Add", apiWrapper(apiHandlerPostsAdd))
	r.HandleFunc("/api/meme.Feed.GetHeader", apiWrapper(apiHandlerFeedGetHeader))
	r.HandleFunc("/api/meme.Utils.ResolveRoute", apiWrapper(apiHandlerResolveRoute))

	// Other
	r.HandleFunc("/vk-callback", handleVKCallback)
	r.HandleFunc("/logout", handleLogout)

	http.Handle("/", r)
	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
