package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/materkov/meme9/web/pb"
)

type Config struct {
	VKAppSecret string
}

var config Config

func apiHandler(method string, body io.Reader, viewerID int) proto.Message {
	switch method {
	case "meme.Feed/Get":
		req := pb.FeedGetRequest{}
		_ = jsonpb.Unmarshal(body, &req)

		return &pb.FeedGetResponse{
			Renderer: &pb.FeedRenderer{
				Posts: []*pb.FeedRenderer_Post{
					{
						Id:           "3331",
						AuthorUrl:    "/profile/4",
						AuthorId:     "4",
						AuthorAvatar: "https://sun3.43222.userapi.com/s/v1/ig2/FGgcvoXeiJaix4uHo4bx7uS1aLgIhTVVbyUqwqXYmTFwNJJJzkLdXXKOiusyXYdqExevW-VSQVytEQ1l2Q3iOSmD.jpg?size=100x0&quality=96&crop=120,33,601,601&ava=1",
						AuthorName:   "Maks Materkov",
						DateDisplay:  "16 янв 2021 в 11:42",
						Text:         "Пост текст текст текст текст текст текст текст текст текст текст текст",
						ImageUrl:     "https://sun9-48.userapi.com/impg/-nlbpow-jdaXTpOb6isHbSobFZQvpt5GvdDaaA/147nS0j68sg.jpg?size=2560x1897&quality=96&sign=97141104582fb9f6a35162fd3aff78ec&type=album",
					},
					{
						Id:           "521",
						AuthorUrl:    "/profile/5",
						AuthorId:     "5",
						AuthorAvatar: "https://sun3.43222.userapi.com/s/v1/ig2/FGgcvoXeiJaix4uHo4bx7uS1aLgIhTVVbyUqwqXYmTFwNJJJzkLdXXKOiusyXYdqExevW-VSQVytEQ1l2Q3iOSmD.jpg?size=100x0&quality=96&crop=120,33,601,601&ava=1",
						AuthorName:   "Vasya Pupmpkin",
						DateDisplay:  "16 апр 2018 в 15:40",
						Text:         "(Note, if you're a new user of ts-proto and using a modern TS setup with esModuleInterop, you need to also pass that as a ts_proto_opt.)",
						ImageUrl:     "https://sun3.43222.userapi.com/impg/ZlHHJIwh9h-Jx8tKkfDB3O5A_XTJgbg0bq_4AQ/crxXzinhDrk.jpg?size=1920x1440&quality=96&sign=c79c25cb7ad78d5ada79408d363593e9&type=album",
					},
				},
			},
		}

	case "meme.Profile/Get":
		req := pb.ProfileGetRequest{}
		_ = jsonpb.Unmarshal(body, &req)

		return &pb.ProfileGetResponse{
			Renderer: &pb.ProfileRenderer{
				Id:     req.Id,
				Name:   fmt.Sprintf("Maks Materkov #%s", req.Id),
				Avatar: "https://sun3.43222.userapi.com/s/v1/ig2/FGgcvoXeiJaix4uHo4bx7uS1aLgIhTVVbyUqwqXYmTFwNJJJzkLdXXKOiusyXYdqExevW-VSQVytEQ1l2Q3iOSmD.jpg?size=100x0&quality=96&crop=120,33,601,601&ava=1",
			},
		}

	case "meme.Feed/GetHeader":
		return &pb.FeedGetHeaderResponse{
			Renderer: &pb.HeaderRenderer{
				MainUrl:    "/",
				UserName:   fmt.Sprintf("User %d", viewerID),
				UserAvatar: "https://sun3.43222.userapi.com/s/v1/ig2/FGgcvoXeiJaix4uHo4bx7uS1aLgIhTVVbyUqwqXYmTFwNJJJzkLdXXKOiusyXYdqExevW-VSQVytEQ1l2Q3iOSmD.jpg?size=100x0&quality=96&crop=120,33,601,601&ava=1",
			},
		}

	case "meme.Posts/Add":
		req := pb.PostsAddRequest{}
		_ = jsonpb.Unmarshal(body, &req)

		return &pb.PostsAddResponse{
			PostUrl: "/profile/5",
		}

	case "meme.Auth/LoginPage":
		requestScheme := "http"
		requestHost := "localhost:8000"
		vkAppID := 7260220
		redirectURL := url.QueryEscape(fmt.Sprintf("%s://%s/vk-callback", requestScheme, requestHost))
		vkURL := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", vkAppID, redirectURL)

		return &pb.LoginPageResponse{
			Renderer: &pb.LoginPageRenderer{
				AuthUrl: vkURL,
				Text:    "Войти через ВК",
			},
		}

	default:
		return nil
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
	resp := apiHandler(method, r.Body, viewerID)

	m := jsonpb.Marshaler{}
	_ = m.Marshal(w, resp)
}

type routeResp struct {
	Component string          `json:"component"`
	Data      json.RawMessage `json:"data"`
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	url := r.URL.Query().Get("url")

	enc := json.NewEncoder(w)
	var component pb.Renderers
	var data proto.Message

	if url == "/" {
		component = pb.Renderers_FEED
		data = apiHandler("meme.Feed/Get", strings.NewReader(""), 0)
	}

	if url == "/login" {
		component = pb.Renderers_LOGIN
		data = apiHandler("meme.Auth/LoginPage", strings.NewReader(""), 0)
	}

	match, _ := regexp.MatchString(`^/profile/(\d+)$`, url)
	if match {
		component = pb.Renderers_PROFILE
		req := &pb.ProfileGetRequest{
			Id: strings.TrimPrefix(url, "/profile/"),
		}
		m := jsonpb.Marshaler{}
		reqStr, _ := m.MarshalToString(req)
		data = apiHandler("meme.Profile/Get", strings.NewReader(reqStr), 0)
	}

	m := jsonpb.Marshaler{}
	dataStr, _ := m.MarshalToString(data)

	_ = enc.Encode(routeResp{
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

	http.HandleFunc("/api/", handler)
	http.HandleFunc("/router", handleRoute)
	http.HandleFunc("/vk-callback", handleVKCallback)

	http.ListenAndServe("127.0.0.1:8000", nil)
}
