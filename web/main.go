package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/materkov/meme9/web/pb"
)

func apiHandler(method string, body io.Reader) proto.Message {
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

	case "meme.Posts/Add":
		req := pb.PostsAddRequest{}
		_ = jsonpb.Unmarshal(body, &req)

		return &pb.PostsAddResponse{
			PostUrl: "/profile/5",
		}

	default:
		return nil
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	method := strings.TrimPrefix(r.URL.Path, "/api/")
	resp := apiHandler(method, r.Body)

	m := jsonpb.Marshaler{}
	_ = m.Marshal(w, resp)
}

type routeResp struct {
	Component string          `json:"component"`
	Data      json.RawMessage `json:"data"`
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	url := r.URL.Query().Get("url")

	enc := json.NewEncoder(w)
	var component pb.Renderers
	var data proto.Message

	if url == "/" {
		component = pb.Renderers_FEED
		data = apiHandler("meme.Feed/Get", strings.NewReader(""))
	}

	match, _ := regexp.MatchString(`^/profile/(\d+)$`, url)
	if match {
		component = pb.Renderers_PROFILE
		req := &pb.ProfileGetRequest{
			Id: strings.TrimPrefix(url, "/profile/"),
		}
		m := jsonpb.Marshaler{}
		reqStr, _ := m.MarshalToString(req)
		data = apiHandler("meme.Profile/Get", strings.NewReader(reqStr))
	}

	m := jsonpb.Marshaler{}
	dataStr, _ := m.MarshalToString(data)

	_ = enc.Encode(routeResp{
		Component: pb.Renderers_name[int32(component)],
		Data:      json.RawMessage(dataStr),
	})
}

func main() {
	http.HandleFunc("/api/", handler)
	http.HandleFunc("/router", handleRoute)

	http.ListenAndServe("127.0.0.1:8000", nil)
}
