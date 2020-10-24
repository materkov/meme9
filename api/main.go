package api

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

func writeResponse(w http.ResponseWriter, resp proto.Message) {
	m := jsonpb.Marshaler{}
	_ = m.Marshal(w, resp)
}

func apiUserPage(viewer *api.Viewer, req *login.UserPageRequest) *login.AnyRenderer {
	return &login.AnyRenderer{Renderer: &login.AnyRenderer_UserPageRenderer{
		UserPageRenderer: &login.UserPageRenderer{
			Id:            req.UserId,
			LastPostId:    "2",
			CurrentUserId: strconv.Itoa(viewer.UserID),
			Name:          req.UserId + " - name",
		},
	}}
}

func apiPostPage(viewer *api.Viewer, req *login.PostPageRequest) *login.AnyRenderer {
	return &login.AnyRenderer{Renderer: &login.AnyRenderer_PostPageRenderer{
		PostPageRenderer: &login.PostPageRenderer{
			Id:            req.PostId,
			Text:          "bla bla bla - " + req.PostId,
			UserId:        "1",
			CurrentUserId: strconv.Itoa(viewer.UserID),
		},
	}}
}

func apiUserConvertUrl(url string) *login.UserPageRequest {
	return &login.UserPageRequest{
		UserId: url[7:],
	}
}

func apiPostConvertUrl(url string) *login.PostPageRequest {
	return &login.PostPageRequest{
		PostId: url[7:],
	}
}

func resolveRoute(url string) resolvedRoute {
	if match, _ := regexp.MatchString(`^/users/([0-9]+)`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_UserPageRequest{
				UserPageRequest: apiUserConvertUrl(url),
			}},
			js: []string{
				"/static/React.js",
				"/static/UserPage.js",
				"/static/Global.js",
			},
		}
	}

	if match, _ := regexp.MatchString(`^/posts/([0-9]+)`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_PostPageRequest{
				PostPageRequest: apiPostConvertUrl(url),
			}},
			js: []string{
				"/static/React.js",
				"/static/PostPage.js",
				"/static/Global.js",
			},
		}
	}

	if match, _ := regexp.MatchString(`^/login`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_LoginPageRequest{}},
			js: []string{
				"/static/React.js",
				"/static/LoginPage.js",
				"/static/Global.js",
			},
		}
	}

	return resolvedRoute{}
}

type resolvedRoute struct {
	apiRequest *login.AnyRequest
	js         []string
}

type Main struct {
	store *store.Store

	loginPage *handlers.LoginPage
}

func (m *Main) apiRequest(viewer *api.Viewer, req *login.AnyRequest) *login.AnyRenderer {
	switch req := req.GetRequest().(type) {
	case *login.AnyRequest_UserPageRequest:
		return apiUserPage(viewer, req.UserPageRequest)
	case *login.AnyRequest_PostPageRequest:
		return apiPostPage(viewer, req.PostPageRequest)
	case *login.AnyRequest_LoginPageRequest:
		return m.loginPage.Handle(viewer, req.LoginPageRequest)
	default:
		return nil
	}
}

func (m *Main) Main() {
	redisClient := redis.NewClient(&redis.Options{})
	m.store = store.NewStore(redisClient)
	authMiddleware := &AuthMiddleware{store: m.store}
	m.loginPage = &handlers.LoginPage{Store: m.store}
	loginApi := &handlers.LoginApi{Store: m.store}

	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		resolvedRoute := resolveRoute(r.URL.Path)
		viewer, _ := r.Context().Value("viewer").(*api.Viewer)

		resp := m.apiRequest(viewer, resolvedRoute.apiRequest)

		page := HTMLPage{
			Request:   resolvedRoute.apiRequest,
			Data:      resp,
			JsBundles: resolvedRoute.js,
			ApiKey:    "access-key",
		}
		_, _ = w.Write([]byte(page.render()))
	}

	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		viewer, _ := r.Context().Value("viewer").(*api.Viewer)

		req := &login.AnyRequest{}
		_ = jsonpb.Unmarshal(r.Body, req)

		resp := m.apiRequest(viewer, req)
		writeResponse(w, resp)
	}

	http.HandleFunc("/api/login", loginApi.ServeHTTP)
	http.HandleFunc("/", authMiddleware.Do(mainHandler))
	http.HandleFunc("/api", authMiddleware.Do(apiHandler))
	http.HandleFunc("/resolve-route", func(w http.ResponseWriter, r *http.Request) {
		req := login.ResolveRouteRequest{}
		_ = jsonpb.Unmarshal(r.Body, &req)

		route := resolveRoute(req.Url)
		writeResponse(w, &login.ResolveRouteResponse{
			Request: route.apiRequest,
			Js:      route.js,
		})
	})

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
