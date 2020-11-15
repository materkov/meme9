package handlers

import (
	"github.com/gogo/protobuf/jsonpb"
	"github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/pkg/api"
	"github.com/materkov/meme9/api/pkg/router"
)

type ResolveRoute struct {
}

func (p *ResolveRoute) Handle(viewer *api.Viewer, req *pb.ResolveRouteRequest) (*pb.ResolveRouteResponse, error) {
	route := router.Resolve(req.Url)

	m := jsonpb.Marshaler{}
	apiRequestStr, _ := m.MarshalToString(route.ApiArgs)

	return &pb.ResolveRouteResponse{
		Js:            route.Js,
		RootComponent: route.RootComponent,
		ApiMethod:     route.ApiMethod,
		ApiRequest:    apiRequestStr,
	}, nil
}
