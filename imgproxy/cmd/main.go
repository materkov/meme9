package main

import (
	"github.com/materkov/meme9/imgproxy"
	"github.com/materkov/meme9/imgproxy/pb/github.com/materkov/meme9/api"
	"net/http"
)

func main() {
	srv := api.NewImageProxyServer(&imgproxy.Service{})
	http.Handle(srv.PathPrefix(), srv)

	_ = http.ListenAndServe("127.0.0.1:8003", nil)
}
