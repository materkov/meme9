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

type apiHandler func(body io.Reader, viewer *Viewer) proto.Message

func apiHandlerFeedGetHeader(body io.Reader, viewer *Viewer) proto.Message {
	headerRenderer := pb.HeaderRenderer{
		MainUrl:   "/",
		LogoutUrl: "http://localhost:8000/logout",
	}

	if viewer.UserID != 0 {
		users, _ := store.GetUsers([]int{viewer.UserID})
		user := users[0]

		headerRenderer.IsAuthorized = true
		headerRenderer.UserAvatar = user.VkAvatar
		headerRenderer.UserName = user.Name
	}

	return &pb.FeedGetHeaderResponse{Renderer: &headerRenderer}
}

func apiHandlerPostsAdd(body io.Reader, viewer *Viewer) proto.Message {
	req := pb.PostsAddRequest{}
	_ = jsonpb.Unmarshal(body, &req)

	post := Post{
		UserID: viewer.UserID,
		Date:   int(time.Now().Unix()),
		Text:   req.Text,
	}

	_ = store.AddPost(&post)

	return &pb.PostsAddResponse{
		PostUrl: fmt.Sprintf("/profile/%d", viewer.UserID),
	}
}

type routeHandler func(url string) *pb.UniversalRenderer

func handleIndex(_ string) *pb.UniversalRenderer {
	postIds, err := store.GetFeed()
	log.Printf("%v %s", postIds, err)

	posts, err := store.GetPosts(postIds)
	log.Printf("%v %s", posts, err)

	wrappedPosts := convertPosts(posts)

	return &pb.UniversalRenderer{
		Renderer: &pb.UniversalRenderer_FeedRenderer{FeedRenderer: &pb.FeedRenderer{
			Posts: wrappedPosts,
		}},
	}
}

func handleLogin(_ string) *pb.UniversalRenderer {
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
	}}
}

func handleProfile(url string) *pb.UniversalRenderer {
	req := &pb.ProfileGetRequest{
		Id: strings.TrimPrefix(url, "/users/"),
	}

	userID, _ := strconv.Atoi(req.Id)
	postIds, _ := store.GetPostsByUser(userID)
	posts, _ := store.GetPosts(postIds)
	wrappedPosts := convertPosts(posts)

	users, _ := store.GetUsers([]int{userID})
	user := users[0]

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_ProfileRenderer{ProfileRenderer: &pb.ProfileRenderer{
		Id:     strconv.Itoa(user.ID),
		Name:   user.Name,
		Avatar: user.VkAvatar,
		Posts:  wrappedPosts,
	}}}
}

func handlePostPage(url string) *pb.UniversalRenderer {
	postIDStr := strings.TrimPrefix(url, "/posts/")
	postID, _ := strconv.Atoi(postIDStr)

	posts, _ := store.GetPosts([]int{postID})
	wrappedPosts := convertPosts(posts)

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_PostRenderer{PostRenderer: &pb.PostRenderer{
		Post: wrappedPosts[0],
	}}}
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

	users, _ := store.GetUsers(userIds)

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

	users, _ := store.GetUsers([]int{userID})
	user := users[0]

	vkName, vkAvatar, err := fetchVkData(body.UserID, body.AccessToken)
	if err == nil {
		user.Name = vkName
		user.VkAvatar = vkAvatar
		_ = store.UpdateNameAvatar(user)
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

		viewer := Viewer{}
		accessCookie, err := r.Cookie("access_token")
		if err == nil && accessCookie.Value != "" {
			token, err := store.GetToken(accessCookie.Value)
			if err == nil {
				viewer.Token = token
				viewer.UserID = token.UserID
			}
		}

		resp := next(r.Body, &viewer)

		m := jsonpb.Marshaler{}
		_ = m.Marshal(w, resp)
	}
}

func routerWrapper(next routeHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")

		data := next(r.URL.RequestURI())
		m := jsonpb.Marshaler{}
		dataStr, _ := m.MarshalToString(data)

		_ = json.NewEncoder(w).Encode(struct {
			//Component string          `json:"component"`
			Data json.RawMessage `json:"data"`
		}{
			//Component: pb.Renderers_name[int32(component)],
			Data: json.RawMessage(dataStr),
		})
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

	// Router
	r.HandleFunc("/", routerWrapper(handleIndex))
	r.HandleFunc("/login", routerWrapper(handleLogin))
	r.HandleFunc("/users/{id}", routerWrapper(handleProfile))
	r.HandleFunc("/posts/{id}", routerWrapper(handlePostPage))

	// Other
	r.HandleFunc("/vk-callback", handleVKCallback)
	r.HandleFunc("/logout", handleLogout)

	http.Handle("/", r)
	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
