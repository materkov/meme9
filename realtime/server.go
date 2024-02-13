package realtime

import (
	"context"
	"github.com/materkov/meme9/realtime/pb/github.com/materkov/meme9/api"
)

type Server struct{}

func (s *Server) GetEvents(ctx context.Context, req *api.GetEventsReq) (*api.GetEventsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) SendEvent(ctx context.Context, req *api.SendEventReq) (*api.Void, error) {
	//TODO implement me
	panic("implement me")
}
