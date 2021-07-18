package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/materkov/meme9/web/pb"
	"github.com/materkov/meme9/web/store"
)

type Relations struct {
	App *App
}

func (r *Relations) Follow(ctx context.Context, req *pb.RelationsFollowRequest) (*pb.RelationsFollowResponse, error) {
	viewer := GetViewerFromContext(ctx)
	if viewer.UserID == 0 {
		return nil, fmt.Errorf("need auth")
	}

	requestedID, _ := strconv.Atoi(req.UserId)
	if requestedID == viewer.UserID {
		return nil, fmt.Errorf("cannot follow yourself")
	}

	obj, err := r.App.Store.ObjGet(ctx, requestedID)
	if err != nil {
	    return nil, fmt.Errorf("error getting user object: %w", err)
	} else if obj == nil || obj.User == nil {
		return nil, fmt.Errorf("user not exists")
	}

	assocs, err := r.App.Store.AssocRange(ctx, viewer.UserID, store.Assoc_Following, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed getting following ids: %w", err)
	}

	followingIds := make([]int, len(assocs))
	for i, assoc := range assocs {
		followingIds[i] = assoc.Following.ID2
	}

	for _, userID := range followingIds {
		if requestedID == userID {
			return &pb.RelationsFollowResponse{}, nil
		}
	}

	err = r.App.Store.AssocAdd(viewer.UserID, requestedID, store.Assoc_Following, &store.StoredAssoc{Following: &store.Following{
		ID1:  viewer.UserID,
		ID2:  requestedID,
		Type: store.Assoc_Liked,
	}})
	if err != nil {
		return nil, fmt.Errorf("failed saving assoc: %w", err)
	}

	return &pb.RelationsFollowResponse{}, nil
}

func (r *Relations) Unfollow(ctx context.Context, req *pb.RelationsUnfollowRequest) (*pb.RelationsUnfollowResponse, error) {
	viewer := GetViewerFromContext(ctx)
	if viewer.UserID == 0 {
		return nil, fmt.Errorf("need auth")
	}

	requestedID, _ := strconv.Atoi(req.UserId)

	err := r.App.Store.AssocDelete(viewer.UserID, requestedID, store.Assoc_Following)
	if err != nil {
		return nil, fmt.Errorf("failed deleting assoc: %w", err)
	}

	return &pb.RelationsUnfollowResponse{}, nil
}

func (r *Relations) Check(ctx context.Context, req *pb.RelationsCheckRequest) (*pb.RelationsCheckResponse, error) {
	viewer := GetViewerFromContext(ctx)
	if viewer.UserID == 0 {
		return nil, fmt.Errorf("need auth")
	}

	requestedID, _ := strconv.Atoi(req.UserId)

	assoc, err := r.App.Store.AssocGet(ctx, viewer.UserID, store.Assoc_Following, requestedID)
	if err != nil {
		return nil, fmt.Errorf("error getting assoc: %w", err)
	}

	return &pb.RelationsCheckResponse{
		IsFollowing: assoc != nil,
	}, nil
}
