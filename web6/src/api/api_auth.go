package api

import (
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

	token := pkg.AuthToken{UserID: userID}
	return &AuthResp{
		Token:    token.ToString(),
		UserID:   strconv.Itoa(userID),
		UserName: user.Name,
	}, nil
}
