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
		writeErrorCode(w, "invalid_request_body", "")
		return
	}

	var registerReq RegisterReq
	err = json.Unmarshal(body, &registerReq)
	if err != nil {
		writeErrorCode(w, "invalid_json", "")
		return
	}

	if registerReq.Username == "" {
		writeErrorCode(w, "username_required", "")
		return
	}
	if registerReq.Password == "" {
		writeErrorCode(w, "password_required", "")
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	// Create user
	user, err := a.users.Create(r.Context(), users.User{
		Username:     registerReq.Username,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	})
	if err != nil {
		if errors.Is(err, users.ErrUsernameExists) {
			writeErrorCode(w, "username_exists", "")
			return
		}
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	// Generate token
	tokenValue, err := generateToken()
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	// Store token
	_, err = a.tokens.Create(r.Context(), tokens.Token{
		Token:     tokenValue,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(LoginResp{
		Token:    tokenValue,
		UserID:   user.ID,
		Username: user.Username,
	})
}
