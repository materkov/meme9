package server

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/materkov/meme9/api/handlers"
	"github.com/materkov/meme9/api/handlers/web"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/pkg/csrf"
	"github.com/materkov/meme9/api/pkg/router"
	"github.com/materkov/meme9/api/store"
)

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

	if err != nil {
		serialized, _ := m.MarshalToString(wrapError(err))
		return `{"ok":false,"data":null,"error":` + serialized + `}`
	} else {
		serialized, _ := m.MarshalToString(resp)
		return `{"ok":true,"data":` + serialized + `}`
	}
}

type Main struct {
	Store  *store.Store
	Config *pkg.Config
}

func (m *Main) Run() {
	authMiddleware := &authMiddleware{store: m.Store}
	csrfMiddleware := &csrfMiddleware{Config: m.Config}
	logoutHandler := &web.Logout{}
	vkCallbackApi := &web.VKCallback{Store: m.Store, Config: m.Config}

	apiHandlers := handlers.NewHandlers(m.Store, m.Config)

	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		route := router.Resolve(r.URL.Path)
		viewer, _ := r.Context().Value(viewerCtxKey).(*api.Viewer)

		// TODO think about this
		encoder := jsonpb.Marshaler{}
		argsStr, _ := encoder.MarshalToString(route.ApiArgs)

		resp, err := apiHandlers.Call(viewer, route.ApiMethod, argsStr)
		initResponse := serializeResponse(resp, err)

		CSRFToken := ""
		if viewer.User != nil {
			CSRFToken = csrf.GenerateToken(m.Config.CSRFKey, viewer.User.ID)
		}

		page := HTMLPage{
			ApiMethod:     route.ApiMethod,
			ApiRequest:    route.ApiArgs,
			ApiResponse:   initResponse,
			JsBundles:     route.Js,
			CSRFToken:     CSRFToken,
			RootComponent: route.RootComponent,
		}
		_, _ = w.Write([]byte(page.Render()))
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

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}
