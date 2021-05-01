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
	return &pb.FeedGetHeaderResponse{
		Renderer: &pb.HeaderRenderer{
			IsAuthorized: viewer.UserID != 0,
			MainUrl:      "/",
			UserName:     fmt.Sprintf("User %d", viewer.UserID),
			UserAvatar:   "https://sun3.43222.userapi.com/s/v1/ig2/FGgcvoXeiJaix4uHo4bx7uS1aLgIhTVVbyUqwqXYmTFwNJJJzkLdXXKOiusyXYdqExevW-VSQVytEQ1l2Q3iOSmD.jpg?size=100x0&quality=96&crop=120,33,601,601&ava=1",
			LogoutUrl:    "http://localhost:8000/logout",
		},
	}
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

	postsMap := map[int]*Post{}
	for _, post := range posts {
		postsMap[post.ID] = post
	}

	postsWrapped := make([]*pb.Post, 0)
	for _, postID := range postIds {
		post := postsMap[postID]
		if post == nil {
			continue
		}

		dateDisplay := time.Unix(int64(post.Date), 0).Format("2 Jan 2006 15:04")

		postsWrapped = append(postsWrapped, &pb.Post{
			Id:           strconv.Itoa(post.ID),
			AuthorUrl:    fmt.Sprintf("/profile/%d", post.ID),
			AuthorId:     strconv.Itoa(post.UserID),
			AuthorAvatar: "https://sun3.43222.userapi.com/s/v1/ig2/FGgcvoXeiJaix4uHo4bx7uS1aLgIhTVVbyUqwqXYmTFwNJJJzkLdXXKOiusyXYdqExevW-VSQVytEQ1l2Q3iOSmD.jpg?size=100x0&quality=96&crop=120,33,601,601&ava=1",
			AuthorName:   "Maks Materkov",
			DateDisplay:  dateDisplay,
			Text:         post.Text,
			ImageUrl:     "https://sun9-48.userapi.com/impg/-nlbpow-jdaXTpOb6isHbSobFZQvpt5GvdDaaA/147nS0j68sg.jpg?size=2560x1897&quality=96&sign=97141104582fb9f6a35162fd3aff78ec&type=album",
		})
	}

	return &pb.UniversalRenderer{
		Renderer: &pb.UniversalRenderer_FeedRenderer{FeedRenderer: &pb.FeedRenderer{
			Posts: postsWrapped,
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
		Id: strings.TrimPrefix(url, "/profile/"),
	}

	userID, _ := strconv.Atoi(req.Id)
	postIds, _ := store.GetPostsByUser(userID)
	posts, _ := store.GetPosts(postIds)
	postsMap := map[int]*Post{}
	for _, post := range posts {
		postsMap[post.ID] = post
	}

	var postsWrapped []*pb.Post
	for _, postID := range postIds {
		post := postsMap[postID]
		if post == nil {
			continue
		}

		postsWrapped = append(postsWrapped, &pb.Post{
			Id:           strconv.Itoa(post.ID),
			AuthorId:     strconv.Itoa(post.UserID),
			AuthorAvatar: "",
			AuthorName:   "",
			AuthorUrl:    fmt.Sprintf("/profile/%d", post.UserID),
			DateDisplay:  "qwer",
			Text:         post.Text,
			ImageUrl:     "",
		})
	}

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_ProfileRenderer{ProfileRenderer: &pb.ProfileRenderer{
		Id:     req.Id,
		Name:   fmt.Sprintf("Maks Materkov #%s", req.Id),
		Avatar: "https://sun3.43222.userapi.com/s/v1/ig2/FGgcvoXeiJaix4uHo4bx7uS1aLgIhTVVbyUqwqXYmTFwNJJJzkLdXXKOiusyXYdqExevW-VSQVytEQ1l2Q3iOSmD.jpg?size=100x0&quality=96&crop=120,33,601,601&ava=1",
		Posts:  postsWrapped,
	}}}
}

func handlePostPage(url string) *pb.UniversalRenderer {
	postIDStr := strings.TrimPrefix(url, "/posts/")
	postID, _ := strconv.Atoi(postIDStr)

	posts, _ := store.GetPosts([]int{postID})
	post := posts[0]

	postWrapped := &pb.Post{
		Id:           strconv.Itoa(post.ID),
		AuthorId:     strconv.Itoa(post.UserID),
		AuthorAvatar: "",
		AuthorName:   "",
		AuthorUrl:    fmt.Sprintf("/profile/%d", post.UserID),
		DateDisplay:  "qwer",
		Text:         post.Text,
		ImageUrl:     "",
	}

	return &pb.UniversalRenderer{Renderer: &pb.UniversalRenderer_PostRenderer{PostRenderer: &pb.PostRenderer{
		Post: postWrapped,
	}}}
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
	r.HandleFunc("/profile/{id}", routerWrapper(handleProfile))
	r.HandleFunc("/posts/{id}", routerWrapper(handlePostPage))

	// Other
	r.HandleFunc("/vk-callback", handleVKCallback)
	r.HandleFunc("/logout", handleLogout)

	http.Handle("/", r)
	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
