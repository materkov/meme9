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
	"github.com/materkov/meme9/web7/api"
)

// RegisterHandler handles user registration requests
type RegisterHandler struct {
	api *api.API
}

// NewRegisterHandler creates a new register handler
func NewRegisterHandler(api *api.API) *RegisterHandler {
	return &RegisterHandler{
		api: api,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Handle processes registration requests
func (h *RegisterHandler) Handle(w http.ResponseWriter, req *http.Request) {
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

	user, err := h.api.CreateUser(req.Context(), users.User{
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

	err = h.api.CreateToken(req.Context(), tokens.Token{
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
