package api

import (
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/store"
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

func (*API) authRegister(_ *Viewer, r *AuthEmailReq) (*AuthResp, error) {
	if r.Email == "" {
		return nil, Error("EmptyEmail")
	}
	if len(r.Email) > 100 {
		return nil, Error("EmailTooLong")
	}
	if r.Password == "" {
		return nil, Error("EmptyPassword")
	}

	userID, err := store.GetEdgeByUniqueKey(store.FakeObjEmailAuth, store.EdgeTypeEmailAuth, r.Email)
	if err != nil {
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
	userID, err = store.AddObject(store.ObjTypeUser, user)
	if err != nil {
		return nil, err
	}
	user.ID = userID

	err = store.AddEdge(store.FakeObjEmailAuth, userID, store.EdgeTypeEmailAuth, r.Email)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = pkg.SendTelegramNotify(fmt.Sprintf("Registration: https://meme.mmaks.me/users/%d", userID))
	}()

	token := pkg.AuthToken{UserID: userID}
	return &AuthResp{
		Token:    token.ToString(),
		UserID:   strconv.Itoa(userID),
		UserName: user.Name,
	}, nil
}

func (*API) authLogin(_ *Viewer, r *AuthEmailReq) (*AuthResp, error) {
	if r.Email == "" || r.Password == "" {
		return nil, Error("InvalidCredentials")
	}

	userID, err := store.GetEdgeByUniqueKey(store.FakeObjEmailAuth, store.EdgeTypeEmailAuth, r.Email)
	if err != nil {
		return nil, err
	} else if userID == 0 {
		return nil, Error("InvalidCredentials")
	}

	user, err := store.GetUser(userID)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(r.Password))
	if err != nil {
		return nil, Error("InvalidCredentials")
	}

	go func() {
		_ = pkg.SendTelegramNotify(fmt.Sprintf("Login: https://meme.mmaks.me/users/%d", userID))
	}()

	token := pkg.AuthToken{UserID: userID}
	return &AuthResp{
		Token:    token.ToString(),
		UserID:   strconv.Itoa(userID),
		UserName: user.Name,
	}, nil
}

type AuthVkReq struct {
	Code        string `json:"code"`
	RedirectURL string `json:"redirectUrl"`
}

func (*API) authVk(_ *Viewer, r *AuthVkReq) (*AuthResp, error) {
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

	userID, err := store.GetEdgeByUniqueKey(store.FakeObjVkAuth, store.EdgeTypeVkAuth, strconv.Itoa(vkUserID))
	if err != nil {
		return nil, err
	}

	if userID == 0 {
		userID, err = store.AddObject(store.ObjTypeUser, &User{
			Name: "VK Auth user",
		})
		if err != nil {
			return nil, err
		}

		err = store.AddEdge(store.FakeObjVkAuth, userID, store.EdgeTypeVkAuth, strconv.Itoa(vkUserID))
		if err != nil {
			return nil, err
		}
	} else {
		user, err := store.GetUser(userID)
		if err != nil {
			return nil, err
		}

		user.Name = userName

		// Already authorized
		err = store.UpdateObject(user, user.ID)
		pkg.LogErr(err)
	}

	token := pkg.AuthToken{UserID: userID}

	go func() {
		_ = pkg.SendTelegramNotify(fmt.Sprintf("VK Auth: https://meme.mmaks.me/users/%d", userID))
	}()

	return &AuthResp{
		Token:    token.ToString(),
		UserID:   strconv.Itoa(userID),
		UserName: userName,
	}, nil
}
