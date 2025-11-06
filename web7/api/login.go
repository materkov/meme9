package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/materkov/meme9/web7/adapters/tokens"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (a *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorCode(w, "invalid_request_body", "")
		return
	}

	var loginReq LoginReq
	err = json.Unmarshal(body, &loginReq)
	if err != nil {
		writeErrorCode(w, "invalid_json", "")
		return
	}

	if loginReq.Username == "" {
		writeErrorCode(w, "username_required", "")
		return
	}
	if loginReq.Password == "" {
		writeErrorCode(w, "password_required", "")
		return
	}

	// Find user by username
	user, err := a.users.GetByUsername(r.Context(), loginReq.Username)
	if err != nil {
		writeErrorCode(w, "invalid_credentials", "")
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		writeErrorCode(w, "invalid_credentials", "")
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

	json.NewEncoder(w).Encode(LoginResp{
		Token:    tokenValue,
		UserID:   user.ID,
		Username: user.Username,
	})
}
