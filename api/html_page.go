package api

import (
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	login "github.com/materkov/meme9/api/pb"
)

type HTMLPage struct {
	Request   *login.AnyRequest
	Data      *login.AnyRenderer
	JsBundles []string
	ApiKey    string
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
		window.InitRequest = {{.InitRequest}};
		window.InitData = {{.InitData}};
		window.InitJsBundles = [{{.InitJsBundles}}];
		window.apiKey = "{{.ApiKey}}";
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
	initDataStr, _ := m.MarshalToString(h.Data)
	initRequestStr, _ := m.MarshalToString(h.Request)

	page = strings.Replace(page, "{{.InitJsBundles}}", jsBundles, 1)
	page = strings.Replace(page, "{{.ApiKey}}", h.ApiKey, 1)
	page = strings.Replace(page, "{{.Scripts}}", scriptTags, 1)
	page = strings.Replace(page, "{{.InitData}}", initDataStr, 1)
	page = strings.Replace(page, "{{.InitRequest}}", initRequestStr, 1)

	return page
}
