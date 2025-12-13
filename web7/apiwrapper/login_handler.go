package apiwrapper

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/materkov/meme9/web7/adapters/tokens"
	"github.com/materkov/meme9/web7/api"
)

// LoginHandler handles login requests
type LoginHandler struct {
	*BaseHandler
}

// NewLoginHandler creates a new login handler
func NewLoginHandler(api *api.API) *LoginHandler {
	return &LoginHandler{
		BaseHandler: NewBaseHandler(api),
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
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

// Handle processes login requests
func (h *LoginHandler) Handle(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeErrorCode(w, "invalid_request_body", "")
		return
	}

	var reqBody LoginRequest
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

	user, err := h.api.GetUserByUsername(req.Context(), reqBody.Username)
	if err != nil {
		writeErrorCode(w, "invalid_credentials", "")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(reqBody.Password))
	if err != nil {
		writeErrorCode(w, "invalid_credentials", "")
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

	json.NewEncoder(w).Encode(LoginResponse{
		Token:    tokenValue,
		UserID:   user.ID,
		Username: user.Username,
	})
}
