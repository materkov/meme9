package server

import (
	"context"
	"fmt"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/pkg"
	"github.com/materkov/meme9/api/src/pkg/utils"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/twitchtv/twirp"
	"strconv"
)

type UserServer struct{}

func transformUser(userID int, user *store.User, viewerID int) (*api.User, error) {
	result := &api.User{
		Id: strconv.Itoa(userID),
	}
	if user == nil {
		return result, nil
	}

	result.Name = user.Name
	result.Status = user.Status

	if viewerID != 0 {
		isFollowing, err := store2.GlobalStore.Subs.CheckFollowing(viewerID, []int{userID})
		if err != nil {
			return nil, err
		} else {
			result.IsFollowing = isFollowing[userID]
		}
	}

	return result, nil
}

func (*UserServer) List(ctx context.Context, r *api.UsersListReq) (*api.UsersList, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	userIds := utils.IdsToInts(r.UserIds)
	users, err := store2.GlobalStore.Users.Get(userIds)
	if err != nil {
		return nil, err
	}

	result := make([]*api.User, len(r.UserIds))
	for i, userIdStr := range r.UserIds {
		userId, _ := strconv.Atoi(userIdStr)

		result[i], err = transformUser(userId, users[userId], viewer.UserID)
		pkg.LogErr(err)
	}

	return &api.UsersList{Users: result}, nil
}

func (*UserServer) SetStatus(ctx context.Context, r *api.UsersSetStatus) (*api.Void, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	if viewer.UserID == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "NotAuthorized")
	}
	if len(r.Status) > 100 {
		return nil, twirp.NewError(twirp.InvalidArgument, "StatusTooLong")
	}

	users, err := store2.GlobalStore.Users.Get([]int{viewer.UserID})
	if err != nil {
		return nil, err
	}

	user := users[viewer.UserID]
	if user == nil {
		return nil, fmt.Errorf("cannot find viewer")
	}

	user.Status = r.Status

	err = store2.GlobalStore.Users.Update(user)
	if err != nil {
		return nil, err
	}

	return &api.Void{}, nil
}

func (*UserServer) Follow(ctx context.Context, r *api.UsersFollowReq) (*api.Void, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	if viewer.UserID == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "NotAuthorized")
	}

	targetID, _ := strconv.Atoi(r.TargetId)
	if targetID <= 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "InvalidTarget")
	}

	if r.Action == api.SubscribeAction_UNFOLLOW {
		err := store2.GlobalStore.Subs.Unfollow(viewer.UserID, targetID)
		pkg.LogErr(err)
	} else {
		users, err := store2.GlobalStore.Users.Get([]int{targetID})
		if err != nil {
			return nil, err
		} else if users[targetID] == nil {
			return nil, twirp.NewError(twirp.InvalidArgument, "UserNotFound")
		}

		err = store2.GlobalStore.Subs.Follow(viewer.UserID, targetID)
		if err != nil {
			return nil, err
		}
	}

	return &api.Void{}, nil
}
