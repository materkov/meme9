package api

import (
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

type HTMLPage struct {
	ApiMethod     string
	ApiRequest    proto.Message
	ApiResponse   proto.Message
	JsBundles     []string
	ApiKey        string
	RootComponent string
}

func (h *HTMLPage) render() string {
	const template = `
<!DOCTYPE html>
<html>
<head>
	<link rel="shortcut icon" href="/static/favicon.ico">
	<meta charset="utf-8">
	<title>meme</title>
</head>
<body>
	<div id="root"></div>
	<script>
		window.modules = {};
		window.InitApiMethod = "{{.InitApiMethod}}";
		window.InitApiRequest = {{.InitApiRequest}};
		window.InitApiResponse = {{.InitApiResponse}};
		window.InitJsBundles = [{{.InitJsBundles}}];
		window.InitRootComponent = "{{.InitRootComponent}}";
		window.apiKey = "{{.ApiKey}}";

		console.log("Initial API:", window.InitApiMethod, window.InitApiRequest, window.InitApiResponse);
	</script>
	{{.Scripts}}
</body>
</html>
`
	page := template

	jsBundles := ""
	scriptTags := ""
	for _, jsBundle := range h.JsBundles {
		jsBundles += `"` + jsBundle + `", `
		scriptTags += `<script src="` + jsBundle + `"></script>`
	}

	m := jsonpb.Marshaler{}
	initApiRequestStr, _ := m.MarshalToString(h.ApiResponse)
	initApiArgsStr, _ := m.MarshalToString(h.ApiRequest)

	page = strings.Replace(page, "{{.InitJsBundles}}", jsBundles, 1)
	page = strings.Replace(page, "{{.ApiKey}}", h.ApiKey, 1)
	page = strings.Replace(page, "{{.Scripts}}", scriptTags, 1)
	page = strings.Replace(page, "{{.InitApiMethod}}", h.ApiMethod, 1)
	page = strings.Replace(page, "{{.InitApiResponse}}", initApiRequestStr, 1)
	page = strings.Replace(page, "{{.InitApiRequest}}", initApiArgsStr, 1)
	page = strings.Replace(page, "{{.InitRootComponent}}", h.RootComponent, 1)

	return page
}
