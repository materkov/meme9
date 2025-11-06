package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/adapters/users"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	var registerReq RegisterReq
	err = json.Unmarshal(body, &registerReq)
	if err != nil {
		writeBadRequest(w, "invalid JSON")
		return
	}

	if registerReq.Username == "" || registerReq.Password == "" {
		writeBadRequest(w, "username and password required")
		return
	}

	// Check if username already exists
	_, err = a.users.GetByUsername(r.Context(), registerReq.Username)
	if err == nil {
		writeConflict(w, "username already exists")
		return
	}
	if !errors.Is(err, users.ErrNotFound) {
		writeInternalServerError(w, "database error")
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		writeInternalServerError(w, "failed to hash password")
		return
	}

	// Create user
	user, err := a.users.Create(r.Context(), users.User{
		Username:     registerReq.Username,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	})
	if err != nil {
		writeInternalServerError(w, "failed to create user")
		return
	}

	// Generate token
	tokenValue, err := generateToken()
	if err != nil {
		writeInternalServerError(w, "failed to generate token")
		return
	}

	// Store token
	_, err = a.tokens.Create(r.Context(), tokens.Token{
		Token:     tokenValue,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		writeInternalServerError(w, "failed to store token")
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(LoginResp{
		Token:    tokenValue,
		UserID:   user.ID,
		Username: user.Username,
	})
}
