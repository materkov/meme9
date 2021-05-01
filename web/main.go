package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	"github.com/jmoiron/sqlx"
	"github.com/materkov/meme9/web/pb"
)

type apiHandler func(body io.Reader, viewerID int) proto.Message

var apiRouter = map[string]apiHandler{
	"meme.Feed/GetHeader": apiHandlerFeedGetHeader,
	"meme.Posts/Add":      apiHandlerPostsAdd,
}

func apiHandlerFeedGetHeader(body io.Reader, viewerID int) proto.Message {
	return &pb.FeedGetHeaderResponse{
		Renderer: &pb.HeaderRenderer{
			IsAuthorized: viewerID != 0,
			MainUrl:      "/",
			UserName:     fmt.Sprintf("User %d", viewerID),
			UserAvatar:   "https://sun3.43222.userapi.com/s/v1/ig2/FGgcvoXeiJaix4uHo4bx7uS1aLgIhTVVbyUqwqXYmTFwNJJJzkLdXXKOiusyXYdqExevW-VSQVytEQ1l2Q3iOSmD.jpg?size=100x0&quality=96&crop=120,33,601,601&ava=1",
		},
	}
}

func apiHandlerPostsAdd(body io.Reader, viewerID int) proto.Message {
	req := pb.PostsAddRequest{}
	_ = jsonpb.Unmarshal(body, &req)

	post := Post{
		UserID: viewerID,
		Date:   int(time.Now().Unix()),
		Text:   req.Text,
	}

	_ = store.AddPost(&post)

	return &pb.PostsAddResponse{
		PostUrl: fmt.Sprintf("/profile/%d", viewerID),
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	viewerID := 0
	accessCookie, err := r.Cookie("access_token")
	if err == nil {
		viewerID, _ = strconv.Atoi(accessCookie.Value)
	}

	method := strings.TrimPrefix(r.URL.Path, "/api/")
	var resp proto.Message
	if apiFunc, ok := apiRouter[method]; ok {
		resp = apiFunc(r.Body, viewerID)
	}

	m := jsonpb.Marshaler{}
	_ = m.Marshal(w, resp)
}

type routeHandler func(url string) (pb.Renderers, proto.Message)

var router = map[*regexp.Regexp]routeHandler{
	regexp.MustCompile(`/`):              handleIndex,
	regexp.MustCompile(`/login`):         handleLogin,
	regexp.MustCompile(`/profile/(\d+)`): handleProfile,
}

func handleIndex(_ string) (pb.Renderers, proto.Message) {
	postIds, err := store.GetFeed()
	log.Printf("%v %s", postIds, err)

	posts, err := store.GetPosts(postIds)
	log.Printf("%v %s", posts, err)

	postsMap := map[int]*Post{}
	for _, post := range posts {
		postsMap[post.ID] = post
	}

	postsWrapped := make([]*pb.FeedRenderer_Post, 0)
	for _, postID := range postIds {
		post := postsMap[postID]
		if post == nil {
			continue
		}

		dateDisplay := time.Unix(int64(post.Date), 0).Format("2 Jan 2006 15:04")

		postsWrapped = append(postsWrapped, &pb.FeedRenderer_Post{
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

	return pb.Renderers_FEED, &pb.FeedGetResponse{
		Renderer: &pb.FeedRenderer{
			Posts: postsWrapped,
		},
	}
}

func handleLogin(_ string) (pb.Renderers, proto.Message) {
	requestScheme := "http"
	requestHost := "localhost:8000"
	vkAppID := 7260220
	redirectURL := url.QueryEscape(fmt.Sprintf("%s://%s/vk-callback", requestScheme, requestHost))
	vkURL := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", vkAppID, redirectURL)

	return pb.Renderers_LOGIN, &pb.LoginPageResponse{
		Renderer: &pb.LoginPageRenderer{
			AuthUrl: vkURL,
			Text:    "Войти через ВК",
		},
	}
}

func handleProfile(url string) (pb.Renderers, proto.Message) {
	req := &pb.ProfileGetRequest{
		Id: strings.TrimPrefix(url, "/profile/"),
	}

	return pb.Renderers_PROFILE, &pb.ProfileGetResponse{
		Renderer: &pb.ProfileRenderer{
			Id:     req.Id,
			Name:   fmt.Sprintf("Maks Materkov #%s", req.Id),
			Avatar: "https://sun3.43222.userapi.com/s/v1/ig2/FGgcvoXeiJaix4uHo4bx7uS1aLgIhTVVbyUqwqXYmTFwNJJJzkLdXXKOiusyXYdqExevW-VSQVytEQ1l2Q3iOSmD.jpg?size=100x0&quality=96&crop=120,33,601,601&ava=1",
		},
	}
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")

	urlAddress := r.URL.Query().Get("url")

	var component pb.Renderers
	var data proto.Message

	for route, handler := range router {
		if route.MatchString(urlAddress) {
			component, data = handler(urlAddress)
		}
	}

	m := jsonpb.Marshaler{}
	dataStr, _ := m.MarshalToString(data)

	_ = json.NewEncoder(w).Encode(struct {
		Component string          `json:"component"`
		Data      json.RawMessage `json:"data"`
	}{
		Component: pb.Renderers_name[int32(component)],
		Data:      json.RawMessage(dataStr),
	})
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

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    strconv.Itoa(body.UserID),
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "http://localhost:3000", http.StatusFound)
}

func main() {
	configStr, _ := os.LookupEnv("CONFIG")
	_ = json.Unmarshal([]byte(configStr), &config)

	db, err := sqlx.Open("mysql", "root:root@/meme9")
	if err != nil {
		panic(err)
	}

	store = Store{db: db}

	http.HandleFunc("/api/", handler)
	http.HandleFunc("/router", handleRoute)
	http.HandleFunc("/vk-callback", handleVKCallback)

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
