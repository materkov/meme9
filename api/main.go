package api

import (
	"net/http"
	"regexp"

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

var globalJs = []string{
	"/static/React.js",
	"/static/Global.js",
}

func resolveRoute(url string) resolvedRoute {
	if match, _ := regexp.MatchString(`^/users/([0-9]+)`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_UserPageRequest{
				UserPageRequest: &login.UserPageRequest{
					UserId: url[7:],
				},
			}},
			js: []string{
				"/static/UserPage.js",
			},
		}
	}

	if match, _ := regexp.MatchString(`^/posts/([0-9]+)`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_PostPageRequest{
				PostPageRequest: &login.PostPageRequest{
					PostId: url[7:],
				},
			}},
			js: []string{
				"/static/PostPage.js",
			},
		}
	}

	if match, _ := regexp.MatchString(`^/login`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_LoginPageRequest{}},
			js: []string{
				"/static/LoginPage.js",
			},
		}
	}

	if match, _ := regexp.MatchString(`^/composer`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_ComposerRequest{}},
			js: []string{
				"/static/Composer.js",
			},
		}
	}

	if match, _ := regexp.MatchString(`^/feed`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_GetFeedRequest{}},
			js: []string{
				"/static/Feed.js",
			},
		}
	}

	if match, _ := regexp.MatchString(`^/vk-callback`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_VkCallbackRequest{}},
			js:         []string{},
		}
	}

	if match, _ := regexp.MatchString(`^/$`, url); match {
		return resolvedRoute{
			apiRequest: &login.AnyRequest{Request: &login.AnyRequest_IndexRequest{}},
			js: []string{
				"/static/Index.js",
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
	addPost   *handlers.AddPost
	getFeed   *handlers.GetFeed
	composer  *handlers.Composer
	index     *handlers.Index
	postPage  *handlers.PostPage
	userPage  *handlers.UserPage
}

func (m *Main) apiRequest(viewer *api.Viewer, req *login.AnyRequest) *login.AnyRenderer {
	switch req := req.GetRequest().(type) {
	case *login.AnyRequest_UserPageRequest:
		return m.userPage.Handle(viewer, req.UserPageRequest)
	case *login.AnyRequest_PostPageRequest:
		return m.postPage.Handle(viewer, req.PostPageRequest)
	case *login.AnyRequest_LoginPageRequest:
		return m.loginPage.Handle(viewer, req.LoginPageRequest)
	case *login.AnyRequest_AddPostRequest:
		return m.addPost.Handle(viewer, req.AddPostRequest)
	case *login.AnyRequest_GetFeedRequest:
		return m.getFeed.Handle(viewer, req.GetFeedRequest)
	case *login.AnyRequest_ComposerRequest:
		return m.composer.Handle(viewer, req.ComposerRequest)
	case *login.AnyRequest_IndexRequest:
		return m.index.Handle(viewer, req.IndexRequest)
	default:
		return nil
	}
}

func (m *Main) Main() {
	redisClient := redis.NewClient(&redis.Options{})
	m.store = store.NewStore(redisClient)
	authMiddleware := &AuthMiddleware{store: m.store}
	m.loginPage = &handlers.LoginPage{Store: m.store}
	m.addPost = &handlers.AddPost{Store: m.store}
	m.getFeed = &handlers.GetFeed{Store: m.store}
	m.composer = &handlers.Composer{Store: m.store}
	m.index = &handlers.Index{Store: m.store}
	loginApi := &handlers.LoginApi{Store: m.store}
	logoutApi := &handlers.LogoutApi{}
	vkCallbackApi := &handlers.VKCallback{Store: m.store}

	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		resolvedRoute := resolveRoute(r.URL.Path)
		viewer, _ := r.Context().Value(viewerCtxKey).(*api.Viewer)

		resp := m.apiRequest(viewer, resolvedRoute.apiRequest)

		js := append(resolvedRoute.js, globalJs...)

		page := HTMLPage{
			Request:   resolvedRoute.apiRequest,
			Data:      resp,
			JsBundles: js,
			ApiKey:    "access-key",
		}
		_, _ = w.Write([]byte(page.render()))
	}

	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		viewer, _ := r.Context().Value(viewerCtxKey).(*api.Viewer)

		req := &login.AnyRequest{}
		_ = jsonpb.Unmarshal(r.Body, req)

		resp := m.apiRequest(viewer, req)
		writeResponse(w, resp)
	}

	http.HandleFunc("/vk-callback", vkCallbackApi.Handle)
	http.HandleFunc("/api/login", loginApi.ServeHTTP)
	http.HandleFunc("/api/logout", logoutApi.ServeHTTP)
	http.HandleFunc("/", authMiddleware.Do(mainHandler))
	http.HandleFunc("/api", authMiddleware.Do(apiHandler))
	http.HandleFunc("/resolve-route", func(w http.ResponseWriter, r *http.Request) {
		req := login.ResolveRouteRequest{}
		_ = jsonpb.Unmarshal(r.Body, &req)

		route := resolveRoute(req.Url)
		js := append(route.js, globalJs...)

		writeResponse(w, &login.ResolveRouteResponse{
			Request: route.apiRequest,
			Js:      js,
		})
	})

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
