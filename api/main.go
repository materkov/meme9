package api

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/materkov/meme9/api/api"
	"github.com/materkov/meme9/api/handlers"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
)

func writeResponse(w http.ResponseWriter, resp proto.Message) {
	w.Header().Set("content-type", "application/json")

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
			js: []string{
				"/static/UserPage.js",
			},
			rootComponent: "UserPage",
			apiMethod:     "meme.API/UserPage",
			apiArgs: &pb.UserPageRequest{
				UserId: url[7:],
			},
		}
	}

	if match, _ := regexp.MatchString(`^/posts/([0-9]+)`, url); match {
		return resolvedRoute{
			js: []string{
				"/static/PostPage.js",
			},
			rootComponent: "PostPage",
			apiMethod:     "meme.API/PostPage",
			apiArgs: &pb.PostPageRequest{
				PostId: url[7:],
			},
		}
	}

	if match, _ := regexp.MatchString(`^/login`, url); match {
		return resolvedRoute{
			js: []string{
				"/static/LoginPage.js",
			},
			rootComponent: "LoginPage",
			apiMethod:     "meme.API/LoginPage",
			apiArgs:       &pb.LoginPageRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/composer`, url); match {
		return resolvedRoute{
			js: []string{
				"/static/Composer.js",
			},
			rootComponent: "Composer",
			apiMethod:     "meme.API/Composer",
			apiArgs:       &pb.ComposerRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/feed`, url); match {
		return resolvedRoute{
			js: []string{
				"/static/Feed.js",
			},
			rootComponent: "Feed",
			apiMethod:     "meme.API/GetFeed",
			apiArgs:       &pb.GetFeedRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/vk-callback`, url); match {
		return resolvedRoute{
			js:            []string{},
			rootComponent: "",
			apiMethod:     "meme.API/VKCallback",
			apiArgs:       &pb.VKCallbackRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/$`, url); match {
		return resolvedRoute{
			js: []string{
				"/static/Index.js",
			},
			rootComponent: "Index",
			apiMethod:     "meme.API/Index",
			apiArgs:       &pb.IndexRequest{},
		}
	}

	return resolvedRoute{}
}

type resolvedRoute struct {
	js            []string
	rootComponent string
	apiMethod     string
	apiArgs       proto.Message
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

func (m *Main) wrapError(err error) *pb.ErrorRenderer {
	errorRenderer := &pb.ErrorRenderer{}

	var apiErr *api.Error
	if errors.As(err, &apiErr) {
		errorRenderer.ErrorCode = apiErr.Code
		errorRenderer.DisplayText = apiErr.DisplayText
	} else {
		log.Printf("[ERROR] Internal error: %s", err)

		errorRenderer.ErrorCode = "INTERNAL_ERROR"
		errorRenderer.DisplayText = "Internal error"
	}

	return errorRenderer
}

func (m *Main) apiRequestV2(viewer *api.Viewer, method string, args string) proto.Message {
	switch method {
	case "meme.API/UserPage":
		req := &pb.UserPageRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return m.wrapError(api.NewError("INVALID_REQUEST", "Failed parsing request"))
		}

		resp, err := m.userPage.Handle(viewer, req)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case "meme.API/PostPage":
		req := &pb.PostPageRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return m.wrapError(api.NewError("INVALID_REQUEST", "Failed parsing request"))
		}

		resp, err := m.postPage.Handle(viewer, req)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case "meme.API/LoginPage":
		req := &pb.LoginPageRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return m.wrapError(api.NewError("INVALID_REQUEST", "Failed parsing request"))
		}

		resp, err := m.loginPage.Handle(viewer, req)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case "meme.API/AddPost":
		req := &pb.AddPostRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return m.wrapError(api.NewError("INVALID_REQUEST", "Failed parsing request"))
		}

		resp, err := m.addPost.Handle(viewer, req)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case "meme.API/GetFeed":
		req := &pb.GetFeedRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return m.wrapError(api.NewError("INVALID_REQUEST", "Failed parsing request"))
		}

		resp, err := m.getFeed.Handle(viewer, req)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case "meme.API/Composer":
		req := &pb.ComposerRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return m.wrapError(api.NewError("INVALID_REQUEST", "Failed parsing request"))
		}

		resp, err := m.composer.Handle(viewer, req)
		if err != nil {
			return m.wrapError(err)
		}

		return resp
	case "meme.API/Index":
		req := &pb.IndexRequest{}
		if err := jsonpb.UnmarshalString(args, req); err != nil {
			return m.wrapError(api.NewError("INVALID_REQUEST", "Failed parsing request"))
		}

		resp, err := m.index.Handle(viewer, req)
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
	csrfMiddleware := &CSRFMiddleware{}
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

		// TODO think about this
		encoder := jsonpb.Marshaler{}
		argsStr, _ := encoder.MarshalToString(resolvedRoute.apiArgs)

		resp := m.apiRequestV2(viewer, resolvedRoute.apiMethod, argsStr)

		js := append(resolvedRoute.js, globalJs...)

		CSRFToken := ""
		if viewer.User != nil {
			CSRFToken = api.GenerateCSRFToken(viewer.User.ID)
		}

		page := HTMLPage{
			ApiMethod:     resolvedRoute.apiMethod,
			ApiRequest:    resolvedRoute.apiArgs,
			ApiResponse:   resp,
			JsBundles:     js,
			ApiKey:        "access-key",
			CSRFToken:     CSRFToken,
			RootComponent: resolvedRoute.rootComponent,
		}
		_, _ = w.Write([]byte(page.render()))
	}

	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		viewer, _ := r.Context().Value(viewerCtxKey).(*api.Viewer)
		body, _ := ioutil.ReadAll(r.Body)

		method := strings.TrimPrefix(r.URL.Path, "/api/")

		resp := m.apiRequestV2(viewer, method, string(body))
		writeResponse(w, resp)
	}

	http.HandleFunc("/vk-callback", vkCallbackApi.Handle)
	http.HandleFunc("/api/logout", logoutApi.ServeHTTP)
	http.HandleFunc("/", authMiddleware.Do(mainHandler))
	http.HandleFunc("/api/", authMiddleware.Do(csrfMiddleware.Do(apiHandler)))
	http.HandleFunc("/resolve-route", func(w http.ResponseWriter, r *http.Request) {
		req := pb.ResolveRouteRequest{}
		_ = jsonpb.Unmarshal(r.Body, &req)

		route := resolveRoute(req.Url)
		js := append(route.js, globalJs...)

		m := jsonpb.Marshaler{}
		apiRequestStr, _ := m.MarshalToString(route.apiArgs)

		writeResponse(w, &pb.ResolveRouteResponse{
			Js:            js,
			RootComponent: route.rootComponent,
			ApiMethod:     route.apiMethod,
			ApiRequest:    apiRequestStr,
		})
	})

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
