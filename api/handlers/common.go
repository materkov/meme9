package handlers

import (
	"log"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	login "github.com/materkov/meme9/api/pb"
)

func writeResponse(w http.ResponseWriter, resp proto.Message) {
	m := jsonpb.Marshaler{}
	_ = m.Marshal(w, resp)
}

func writeError(w http.ResponseWriter, err *login.ErrorRenderer) {
	writeResponse(w, &login.AnyRenderer{Renderer: &login.AnyRenderer_ErrorRenderer{
		ErrorRenderer: err,
	}})
}

func writeInternalError(w http.ResponseWriter, err error) {
	log.Printf("[ERROR] Internal error: %s", err)
	writeError(w, &login.ErrorRenderer{DisplayText: "Internal error"})
}

const VKAppID = 7260220
