package api

import (
	"net/http"
	"regexp"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	login "github.com/materkov/meme9/api/pb"
)

func writeResponse(w http.ResponseWriter, resp proto.Message) {
	m := jsonpb.Marshaler{}
	_ = m.Marshal(w, resp)
}

func apiUserPage(req *login.UserPageRequest) *login.AnyRenderer {
	return &login.AnyRenderer{Renderer: &login.AnyRenderer_UserPageRenderer{
		UserPageRenderer: &login.UserPageRenderer{
			Id:            req.UserId,
			LastPostId:    "2",
			CurrentUserId: "-102-13",
			Name:          req.UserId + " - name",
		},
	}}
}

func apiPostPage(req *login.PostPageRequest) *login.AnyRenderer {
	return &login.AnyRenderer{Renderer: &login.AnyRenderer_PostPageRenderer{
		PostPageRenderer: &login.PostPageRenderer{
			Id:            req.PostId,
			Text:          "bla bla bla - " + req.PostId,
			UserId:        "1",
			CurrentUserId: "-102-13",
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

	return resolvedRoute{}
}

type resolvedRoute struct {
	apiRequest *login.AnyRequest
	js         []string
}

func apiRequest(req *login.AnyRequest) *login.AnyRenderer {
	switch request := req.GetRequest().(type) {
	case *login.AnyRequest_UserPageRequest:
		return apiUserPage(request.UserPageRequest)
	case *login.AnyRequest_PostPageRequest:
		return apiPostPage(request.PostPageRequest)
	default:
		return nil
	}
}

func Main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resolvedRoute := resolveRoute(r.URL.Path)

		resp := apiRequest(resolvedRoute.apiRequest)

		page := HTMLPage{
			Request:   resolvedRoute.apiRequest,
			Data:      resp,
			JsBundles: resolvedRoute.js,
			ApiKey:    "access-key",
		}
		_, _ = w.Write([]byte(page.render()))
	})
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		req := &login.AnyRequest{}
		_ = jsonpb.Unmarshal(r.Body, req)

		resp := apiRequest(req)
		writeResponse(w, resp)
	})
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
