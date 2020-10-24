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
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return
	}

	tokenID, err := l.Store.GenerateNodeID()
	if err != nil {
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
}
