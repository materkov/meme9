package api

import (
	"errors"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AuthEmailReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResp struct {
	Token    string `json:"token"`
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
}

func (a *API) authRegister(_ *Viewer, r *AuthEmailReq) (*AuthResp, error) {
	if r.Email == "" {
		return nil, Error("EmptyEmail")
	}
	if len(r.Email) > 100 {
		return nil, Error("EmailTooLong")
	}
	if r.Password == "" {
		return nil, Error("EmptyPassword")
	}

	userID, err := store2.GlobalStore.Unique.Get(store2.UniqueTypeEmail, r.Email)
	if err != nil && !errors.Is(err, store2.ErrNotFound) {
		return nil, err
	} else if userID != 0 {
		return nil, Error("EmailAlreadyRegistered")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := store.User{
		Name:         r.Email,
		PasswordHash: string(passwordHash),
	}
	userID, err = store2.GlobalStore.Nodes.Add(store.ObjTypeUser, user)
	if err != nil {
		return nil, err
	}
	user.ID = userID

	err = store2.GlobalStore.Unique.Add(store2.UniqueTypeEmail, r.Email, userID)
	if err != nil {
		return nil, err
	}

	pkg.SendTelegramNotifyAsync(fmt.Sprintf("Registration: https://meme.mmaks.me/users/%d", userID))

	token, err := pkg.GenerateAuthToken(userID)
	if err != nil {
		return nil, err
	}

	return &AuthResp{
		Token:    token,
		UserID:   strconv.Itoa(userID),
		UserName: user.Name,
	}, nil
}

func (a *API) authLogin(_ *Viewer, r *AuthEmailReq) (*AuthResp, error) {
	if r.Email == "" || r.Password == "" {
		return nil, Error("InvalidCredentials")
	}

	userID, err := store2.GlobalStore.Unique.Get(store2.UniqueTypeEmail, r.Email)
	if errors.Is(err, store2.ErrNotFound) {
		return nil, Error("InvalidCredentials")
	} else if err != nil {
		return nil, err
	}

	user, err := store2.GlobalStore.TypedNodes.GetUser(userID)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(r.Password))
	if err != nil {
		return nil, Error("InvalidCredentials")
	}

	pkg.SendTelegramNotifyAsync(fmt.Sprintf("Login: https://meme.mmaks.me/users/%d", userID))

	token, err := pkg.GenerateAuthToken(userID)
	if err != nil {
		return nil, err
	}

	return &AuthResp{
		Token:    token,
		UserID:   strconv.Itoa(userID),
		UserName: user.Name,
	}, nil
}

type AuthVkReq struct {
	Code        string `json:"code"`
	RedirectURL string `json:"redirectUrl"`
}

func (a *API) authVk(_ *Viewer, r *AuthVkReq) (*AuthResp, error) {
	if r.Code == "" {
		return nil, Error("InvalidCode")
	}

	vkUserID, accessToken, err := pkg.ExchangeCode(r.Code, r.RedirectURL)
	if err != nil {
		return nil, Error("InvalidCode")
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
		userID, err = store2.GlobalStore.Nodes.Add(store.ObjTypeUser, &store.User{
			Name: "VK Auth user",
		})
		if err != nil {
			return nil, err
		}

		err = store2.GlobalStore.Unique.Add(store2.UniqueTypeVKID, strconv.Itoa(vkUserID), userID)
		if err != nil {
			return nil, err
		}
	} else {
		user, err := store2.GlobalStore.TypedNodes.GetUser(userID)
		if err != nil {
			return nil, err
		}

		user.Name = userName

		// Already authorized
		err = store2.GlobalStore.Nodes.Update(user.ID, user)
		pkg.LogErr(err)
	}

	pkg.SendTelegramNotifyAsync(fmt.Sprintf("VK Auth: https://meme.mmaks.me/users/%d", userID))

	token, err := pkg.GenerateAuthToken(userID)
	if err != nil {
		return nil, err
	}

	return &AuthResp{
		Token:    token,
		UserID:   strconv.Itoa(userID),
		UserName: userName,
	}, nil
}
