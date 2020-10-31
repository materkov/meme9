package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	login "github.com/materkov/meme9/api/pb"
	"github.com/materkov/meme9/api/store"
	"golang.org/x/crypto/bcrypt"
)

type LoginApi struct {
	Store *store.Store
}

func (l *LoginApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := login.LoginRequest{}
	err := jsonpb.Unmarshal(r.Body, &req)
	if err != nil {
		return
	}

	nodeID, _ := strconv.Atoi(req.Login)

	user, err := l.Store.GetUser(nodeID)
	if err == store.ErrNodeNotFound {
		writeError(w, &login.ErrorRenderer{DisplayText: "Неправильный логин или пароль"})
		return
	} else if err != nil {
		writeInternalError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		writeError(w, &login.ErrorRenderer{DisplayText: "Неправильный логин или пароль"})
		return
	}

	tokenID, err := l.Store.GenerateNodeID()
	if err != nil {
		writeInternalError(w, err)
		return
	}

	token := store.Token{
		ID:       tokenID,
		IssuedAt: int(time.Now().Unix()),
		UserID:   user.ID,
	}
	token.GenerateRandomToken()

	err = l.Store.AddToken(&token)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.Token,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	writeResponse(w, &login.AnyRenderer{Renderer: &login.AnyRenderer_LoginRenderer{
		LoginRenderer: &login.LoginRenderer{
			HeaderRenderer: &login.HeaderRenderer{
				CurrentUserId:   strconv.Itoa(user.ID),
				CurrentUserName: user.Name,
			},
		},
	}})
}
