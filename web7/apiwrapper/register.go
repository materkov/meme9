package apiwrapper

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

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *Router) registerHandler(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeErrorCode(w, "invalid_request_body", "")
		return
	}

	var reqBody RegisterRequest
	if err := json.Unmarshal(body, &reqBody); err != nil {
		writeErrorCode(w, "invalid_json", "")
		return
	}

	if reqBody.Username == "" {
		writeErrorCode(w, "username_required", "")
		return
	}
	if reqBody.Password == "" {
		writeErrorCode(w, "password_required", "")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	user, err := r.api.CreateUser(req.Context(), users.User{
		Username:     reqBody.Username,
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

	tokenValue, err := generateToken()
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	err = r.api.CreateToken(req.Context(), tokens.Token{
		Token:     tokenValue,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		writeInternalServerError(w, "internal_server_error", "")
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(LoginResponse{
		Token:    tokenValue,
		UserID:   user.ID,
		Username: user.Username,
	})
}
