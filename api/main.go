package api

import (
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers"
	"github.com/materkov/meme9/api/pb"
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
			apiRequest: &pb.AnyRequest{Request: &pb.AnyRequest_UserPageRequest{
				UserPageRequest: &pb.UserPageRequest{
					UserId: url[7:],
				},
			}},
			js: []string{
				"/static/UserPage.js",
			},
			rootComponent: "UserPage",
		}
	}

	if match, _ := regexp.MatchString(`^/posts/([0-9]+)`, url); match {
		return resolvedRoute{
			apiRequest: &pb.AnyRequest{Request: &pb.AnyRequest_PostPageRequest{
				PostPageRequest: &pb.PostPageRequest{
					PostId: url[7:],
				},
			}},
			js: []string{
				"/static/PostPage.js",
			},
			rootComponent: "PostPage",
		}
	}

	if match, _ := regexp.MatchString(`^/login`, url); match {
		return resolvedRoute{
			apiRequest: &pb.AnyRequest{Request: &pb.AnyRequest_LoginPageRequest{}},
			js: []string{
				"/static/LoginPage.js",
			},
			rootComponent: "LoginPage",
		}
	}

	if match, _ := regexp.MatchString(`^/composer`, url); match {
		return resolvedRoute{
			apiRequest: &pb.AnyRequest{Request: &pb.AnyRequest_ComposerRequest{}},
			js: []string{
				"/static/Composer.js",
			},
			rootComponent: "Composer",
		}
	}

	if match, _ := regexp.MatchString(`^/feed`, url); match {
		return resolvedRoute{
			apiRequest: &pb.AnyRequest{Request: &pb.AnyRequest_GetFeedRequest{}},
			js: []string{
				"/static/Feed.js",
			},
			rootComponent: "Feed",
		}
	}

	if match, _ := regexp.MatchString(`^/vk-callback`, url); match {
		return resolvedRoute{
			apiRequest:    &pb.AnyRequest{Request: &pb.AnyRequest_VkCallbackRequest{}},
			js:            []string{},
			rootComponent: "",
		}
	}

	if match, _ := regexp.MatchString(`^/$`, url); match {
		return resolvedRoute{
			apiRequest: &pb.AnyRequest{Request: &pb.AnyRequest_IndexRequest{}},
			js: []string{
				"/static/Index.js",
			},
			rootComponent: "Index",
		}
	}

	return resolvedRoute{}
}

type resolvedRoute struct {
	apiRequest    *pb.AnyRequest
	js            []string
	rootComponent string
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

func (m *Main) wrapError(err error) *pb.AnyRenderer {
	errorRenderer := pb.ErrorRenderer{}

	var apiErr *api.Error
	if errors.As(err, &apiErr) {
		errorRenderer.ErrorCode = apiErr.Code
		errorRenderer.DisplayText = apiErr.DisplayText
	} else {
		log.Printf("[ERROR] Internal error: %s", err)

		errorRenderer.ErrorCode = "INTERNAL_ERROR"
		errorRenderer.DisplayText = "Internal error"
	}

	return &pb.AnyRenderer{Renderer: &pb.AnyRenderer_ErrorRenderer{
		ErrorRenderer: &errorRenderer,
	}}
}

func (m *Main) apiRequest(viewer *api.Viewer, req *pb.AnyRequest) proto.Message {
	switch req := req.GetRequest().(type) {
	case *pb.AnyRequest_UserPageRequest:
		resp, err := m.userPage.Handle(viewer, req.UserPageRequest)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case *pb.AnyRequest_PostPageRequest:
		resp, err := m.postPage.Handle(viewer, req.PostPageRequest)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case *pb.AnyRequest_LoginPageRequest:
		resp, err :=  m.loginPage.Handle(viewer, req.LoginPageRequest)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case *pb.AnyRequest_AddPostRequest:
		resp, err := m.addPost.Handle(viewer, req.AddPostRequest)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case *pb.AnyRequest_GetFeedRequest:
		resp, err := m.getFeed.Handle(viewer, req.GetFeedRequest)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case *pb.AnyRequest_ComposerRequest:
		resp, err := m.composer.Handle(viewer, req.ComposerRequest)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case *pb.AnyRequest_IndexRequest:
		resp, err := m.index.Handle(viewer, req.IndexRequest)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
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
	m.postPage = &handlers.PostPage{Store: m.store}
	m.userPage = &handlers.UserPage{Store: m.store}
	m.getFeed = &handlers.GetFeed{Store: m.store}
	m.composer = &handlers.Composer{Store: m.store}
	m.index = &handlers.Index{Store: m.store}
	logoutApi := &handlers.LogoutApi{}
	vkCallbackApi := &handlers.VKCallback{Store: m.store}

	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		resolvedRoute := resolveRoute(r.URL.Path)
		viewer, _ := r.Context().Value(viewerCtxKey).(*api.Viewer)

		resp := m.apiRequest(viewer, resolvedRoute.apiRequest)

		js := append(resolvedRoute.js, globalJs...)

		page := HTMLPage{
			Request:       resolvedRoute.apiRequest,
			Data:          resp,
			JsBundles:     js,
			ApiKey:        "access-key",
			RootComponent: resolvedRoute.rootComponent,
		}
		_, _ = w.Write([]byte(page.render()))
	}

	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		viewer, _ := r.Context().Value(viewerCtxKey).(*api.Viewer)

		req := &pb.AnyRequest{}
		_ = jsonpb.Unmarshal(r.Body, req)

		resp := m.apiRequest(viewer, req)
		writeResponse(w, resp)
	}

	http.HandleFunc("/vk-callback", vkCallbackApi.Handle)
	http.HandleFunc("/api/logout", logoutApi.ServeHTTP)
	http.HandleFunc("/", authMiddleware.Do(mainHandler))
	http.HandleFunc("/api", authMiddleware.Do(apiHandler))
	http.HandleFunc("/resolve-route", func(w http.ResponseWriter, r *http.Request) {
		req := pb.ResolveRouteRequest{}
		_ = jsonpb.Unmarshal(r.Body, &req)

		route := resolveRoute(req.Url)
		js := append(route.js, globalJs...)

		writeResponse(w, &pb.ResolveRouteResponse{
			Request:       route.apiRequest,
			Js:            js,
			RootComponent: route.rootComponent,
		})
	})

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
