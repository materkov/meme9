package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/pkg"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/twitchtv/twirp"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AuthServer struct{}

func (a *AuthServer) CheckAuth(ctx context.Context, req *api.CheckAuthReq) (*api.AuthResp, error) {
	token := pkg.ParseAuthToken(ctx, req.Token)
	if token == nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "IncorrectToken")
	}

	users, err := store2.GlobalStore.Users.Get([]int{token.UserID})
	if err != nil {
		return nil, err
	}

	user := users[token.UserID]
	if user == nil {
		return nil, twirp.NewErrorf(twirp.Internal, "cannot find user %d from auth token", token.UserID)
	}

	return &api.AuthResp{
		UserId:   strconv.Itoa(user.ID),
		UserName: user.Name,
	}, err
}

// TODO add tracers
func (a *AuthServer) Register(ctx context.Context, r *api.EmailReq) (*api.AuthResp, error) {
	if r.Email == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "EmptyEmail")
	}
	if len(r.Email) > 100 {
		return nil, twirp.NewError(twirp.InvalidArgument, "EmailTooLong")
	}
	if r.Password == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "EmptyPassword")
	}

	userID, err := store2.GlobalStore.Unique.Get(store2.UniqueTypeEmail, r.Email)
	if err != nil && !errors.Is(err, store2.ErrNotFound) {
		return nil, err
	} else if userID != 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "EmailAlreadyRegistered")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := store.User{
		Name:         r.Email,
		PasswordHash: string(passwordHash),
	}
	err = store2.GlobalStore.Users.Add(&user)
	if err != nil {
		return nil, err
	}

	err = store2.GlobalStore.Unique.Add(store2.UniqueTypeEmail, r.Email, user.ID)
	if err != nil {
		return nil, err
	}

	pkg.SendTelegramNotifyAsync(fmt.Sprintf("Registration: https://meme.mmaks.me/users/%d", user.ID))

	token, err := pkg.GenerateAuthToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &api.AuthResp{
		Token:    token,
		UserId:   strconv.Itoa(user.ID),
		UserName: user.Name,
	}, nil
}

func (a *AuthServer) Login(ctx context.Context, r *api.EmailReq) (*api.AuthResp, error) {
	if r.Email == "" || r.Password == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "InvalidCredentials")
	}

	userID, err := store2.GlobalStore.Unique.Get(store2.UniqueTypeEmail, r.Email)
	if errors.Is(err, store2.ErrNotFound) {
		return nil, twirp.NewError(twirp.InvalidArgument, "InvalidCredentials")
	} else if err != nil {
		return nil, err
	}

	users, err := store2.GlobalStore.Users.Get([]int{userID})
	if err != nil {
		return nil, err
	} else if users[userID] == nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "InvalidCredentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(users[userID].PasswordHash), []byte(r.Password))
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "InvalidCredentials")
	}

	pkg.SendTelegramNotifyAsync(fmt.Sprintf("Login: https://meme.mmaks.me/users/%d", userID))

	token, err := pkg.GenerateAuthToken(userID)
	if err != nil {
		return nil, err
	}

	return &api.AuthResp{
		Token:    token,
		UserId:   strconv.Itoa(userID),
		UserName: users[userID].Name,
	}, nil
}

func (a *AuthServer) Vk(ctx context.Context, r *api.VkReq) (*api.AuthResp, error) {
	if r.Code == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "InvalidCode")
	}

	vkUserID, accessToken, err := pkg.ExchangeCode(r.Code, r.RedirectUrl)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "InvalidCode")
	}

	userName, err := pkg.RefreshFromVk(accessToken, vkUserID)
	if err != nil {
		return nil, err
	}

	userID, err := store2.GlobalStore.Unique.Get(store2.UniqueTypeVKID, strconv.Itoa(vkUserID))
	if err != nil && !errors.Is(err, store2.ErrNotFound) {
		return nil, err
	}

	if errors.Is(err, store2.ErrNotFound) {
		user := &store.User{
			Name: "VK Auth user",
		}
		err = store2.GlobalStore.Users.Add(user)
		if err != nil {
			return nil, err
		}
		userID = user.ID

		err = store2.GlobalStore.Unique.Add(store2.UniqueTypeVKID, strconv.Itoa(vkUserID), userID)
		if err != nil {
			return nil, err
		}
	} else {
		users, err := store2.GlobalStore.Users.Get([]int{userID})
		if err != nil {
			return nil, err
		} else if users[userID] == nil {
			return nil, fmt.Errorf("cannot find user")
		}

		user := users[userID]
		user.Name = userName

		// Already authorized
		err = store2.GlobalStore.Users.Update(user)
		pkg.LogErr(err)
	}

	pkg.SendTelegramNotifyAsync(fmt.Sprintf("VK Auth: https://meme.mmaks.me/users/%d", userID))

	token, err := pkg.GenerateAuthToken(userID)
	if err != nil {
		return nil, err
	}

	return &api.AuthResp{
		Token:    token,
		UserId:   strconv.Itoa(userID),
		UserName: userName,
	}, nil
}
