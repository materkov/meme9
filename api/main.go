package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/materkov/meme9/api/handlers"
	"github.com/materkov/meme9/api/handlers/web"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/pkg/config"
	"github.com/materkov/meme9/api/pkg/csrf"
	"github.com/materkov/meme9/api/store"
)

func writeResponse(w http.ResponseWriter, resp proto.Message) {
	w.Header().Set("content-type", "application/json")

	m := jsonpb.Marshaler{}
	_ = m.Marshal(w, resp)
}

func wrapError(err error) *pb.ErrorRenderer {
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

func serializeResponse(resp proto.Message, err error) string {
	m := jsonpb.Marshaler{}

	dataStr := ""
	errorStr := ""
	okStr := ""

	if err != nil {
		okStr = "false"
		errorStr, _ = m.MarshalToString(wrapError(err))
		errorStr = `, "error": ` + errorStr
	} else {
		okStr = "true"
		dataStr, _ = m.MarshalToString(resp)
	}

	return fmt.Sprintf(`{"ok": %s, "data": %s%s}`, okStr, dataStr, errorStr)
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
			apiMethod:     "meme.API.UserPage",
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
			apiMethod:     "meme.API.PostPage",
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
			apiMethod:     "meme.API.LoginPage",
			apiArgs:       &pb.LoginPageRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/composer`, url); match {
		return resolvedRoute{
			js: []string{
				"/static/Composer.js",
			},
			rootComponent: "Composer",
			apiMethod:     "meme.API.Composer",
			apiArgs:       &pb.ComposerRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/feed`, url); match {
		return resolvedRoute{
			js: []string{
				"/static/Feed.js",
			},
			rootComponent: "Feed",
			apiMethod:     "meme.API.GetFeed",
			apiArgs:       &pb.GetFeedRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/vk-callback`, url); match {
		return resolvedRoute{
			js:            []string{},
			rootComponent: "",
			apiMethod:     "meme.API.VKCallback",
			apiArgs:       &pb.VKCallbackRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/$`, url); match {
		return resolvedRoute{
			js: []string{
				"/static/Index.js",
			},
			rootComponent: "Index",
			apiMethod:     "meme.API.Index",
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
	store  *store.Store
	Config *config.Config
}

func (m *Main) Main() {
	redisClient := redis.NewClient(&redis.Options{})
	m.store = store.NewStore(redisClient)
	authMiddleware := &AuthMiddleware{store: m.store}
	csrfMiddleware := &CSRFMiddleware{Config: m.Config}
	logoutHandler := &web.Logout{}
	vkCallbackApi := &web.VKCallback{Store: m.store, Config: m.Config}

	apiHandlers := handlers.NewHandlers(m.store, m.Config)

	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		resolvedRoute := resolveRoute(r.URL.Path)
		viewer, _ := r.Context().Value(viewerCtxKey).(*api.Viewer)

		// TODO think about this
		encoder := jsonpb.Marshaler{}
		argsStr, _ := encoder.MarshalToString(resolvedRoute.apiArgs)

		resp, err := apiHandlers.Call(viewer, resolvedRoute.apiMethod, argsStr)
		initResponse := serializeResponse(resp, err)

		js := append(resolvedRoute.js, globalJs...)

		CSRFToken := ""
		if viewer.User != nil {
			CSRFToken = csrf.GenerateCSRFToken(m.Config.CSRFKey, viewer.User.ID)
		}

		page := HTMLPage{
			ApiMethod:     resolvedRoute.apiMethod,
			ApiRequest:    resolvedRoute.apiArgs,
			ApiResponse:   initResponse,
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

		method := r.URL.Query().Get("method")

		resp, err := apiHandlers.Call(viewer, method, string(body))

		respStr := serializeResponse(resp, err)
		_, _ = w.Write([]byte(respStr))
	}

	http.HandleFunc("/vk-callback", vkCallbackApi.Handle)
	http.HandleFunc("/logout", logoutHandler.ServeHTTP)
	http.HandleFunc("/", authMiddleware.Do(mainHandler))
	http.HandleFunc("/api", authMiddleware.Do(csrfMiddleware.Do(apiHandler)))
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
