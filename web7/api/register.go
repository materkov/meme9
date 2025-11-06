package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	mongodriver "go.mongodb.org/mongo-driver/mongo"

	"github.com/materkov/meme9/web7/adapters/mongo"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	var registerReq RegisterReq
	err = json.Unmarshal(body, &registerReq)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	if registerReq.Username == "" || registerReq.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "username and password required"})
		return
	}

	// Check if username already exists
	_, err = a.mongo.GetUserByUsername(r.Context(), registerReq.Username)
	if err == nil {
		w.WriteHeader(409)
		json.NewEncoder(w).Encode(map[string]string{"error": "username already exists"})
		return
	}
	if !errors.Is(err, mongodriver.ErrNoDocuments) {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to hash password"})
		return
	}

	// Create user
	user, err := a.mongo.CreateUser(r.Context(), mongo.User{
		Username:     registerReq.Username,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create user"})
		return
	}

	// Generate token
	tokenValue, err := generateToken()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate token"})
		return
	}

	// Store token
	_, err = a.mongo.CreateToken(r.Context(), mongo.Token{
		Token:     tokenValue,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to store token"})
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(LoginResp{
		Token:    tokenValue,
		UserID:   user.ID,
		Username: user.Username,
	})
}
