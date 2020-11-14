package api

import (
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

type HTMLPage struct {
	ApiMethod     string
	ApiRequest    proto.Message
	ApiResponse   string
	JsBundles     []string
	ApiKey        string
	RootComponent string
	CSRFToken     string
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
		window.CSRFToken = "{{.CSRFToken}}";

		console.log("Initial API", {method: window.InitApiMethod, request: window.InitApiRequest, response: window.InitApiResponse});
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
	initApiArgsStr, _ := m.MarshalToString(h.ApiRequest)

	page = strings.Replace(page, "{{.InitJsBundles}}", jsBundles, 1)
	page = strings.Replace(page, "{{.ApiKey}}", h.ApiKey, 1)
	page = strings.Replace(page, "{{.Scripts}}", scriptTags, 1)
	page = strings.Replace(page, "{{.InitApiMethod}}", h.ApiMethod, 1)
	page = strings.Replace(page, "{{.InitApiResponse}}", h.ApiResponse, 1)
	page = strings.Replace(page, "{{.InitApiRequest}}", initApiArgsStr, 1)
	page = strings.Replace(page, "{{.InitRootComponent}}", h.RootComponent, 1)
	page = strings.Replace(page, "{{.CSRFToken}}", h.CSRFToken, 1)

	return page
}
