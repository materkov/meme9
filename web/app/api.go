package app

import (
	"context"
	"fmt"
	"time"

	"github.com/materkov/meme9/web/pb"
	"github.com/materkov/meme9/web/tracer"
	"github.com/materkov/meme9/web/utils"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (a *App) HandleJSONRequest(ctx context.Context, method string, request []byte) ([]byte, error) {
	currentSpan, ok := ctx.Value(utils.RequestIdKey{}).(int)
	if ok {
		childTracer := &tracer.Tracer{
			Started: time.Now(),
			Name:    fmt.Sprintf("API %s", method),
			TraceID: currentSpan,
		}
		defer childTracer.Stop()
	}

	Logf(ctx, "API %s, req: %s", method, request)

	m := protojson.UnmarshalOptions{DiscardUnknown: true}

	var resp proto.Message
	var err error

	switch method {
	case "meme.Feed.GetHeader":
		req := &pb.FeedGetHeaderRequest{}
		if err := m.Unmarshal(request, req); err != nil {
			return nil, fmt.Errorf("failed unmarshaling request")
		}
		resp, err = FeedSrv.GetHeader(ctx, req)
	case "meme.Posts.Add":
		req := &pb.PostsAddRequest{}
		if err := m.Unmarshal(request, req); err != nil {
			return nil, fmt.Errorf("failed unmarshaling request")
		}
		resp, err = postsSrv.Add(ctx, req)
	case "meme.Posts.ToggleLike":
		req := &pb.ToggleLikeRequest{}
		if err := m.Unmarshal(request, req); err != nil {
			return nil, fmt.Errorf("failed unmarshaling request")
		}
		resp, err = postsSrv.ToggleLike(ctx, req)
	case "meme.Posts.AddComment":
		req := &pb.AddCommentRequest{}
		if err := m.Unmarshal(request, req); err != nil {
			return nil, fmt.Errorf("failed unmarshaling request")
		}
		resp, err = postsSrv.AddComment(ctx, req)
	case "meme.Utils.ResolveRoute":
		req := &pb.ResolveRouteRequest{}
		if err := m.Unmarshal(request, req); err != nil {
			return nil, fmt.Errorf("failed unmarshaling request")
		}
		resp, err = UtilsSrv.ResolveRoute(ctx, req)
	case "meme.Relations.Follow":
		req := &pb.RelationsFollowRequest{}
		if err := m.Unmarshal(request, req); err != nil {
			return nil, fmt.Errorf("failed unmarshaling request")
		}
		resp, err = relationsSrv.Follow(ctx, req)
	case "meme.Relations.Unfollow":
		req := &pb.RelationsUnfollowRequest{}
		if err := m.Unmarshal(request, req); err != nil {
			return nil, fmt.Errorf("failed unmarshaling request")
		}
		resp, err = relationsSrv.Unfollow(ctx, req)
	default:
		return nil, fmt.Errorf("unknown method")
	}

	if err != nil {
		Logf(ctx, "API error: %s", err)
		return nil, err
	}

	marshaller := &protojson.MarshalOptions{}
	respBytes, _ := marshaller.Marshal(resp)
	Logf(ctx, "API response: %s", respBytes)

	return respBytes, nil
}
