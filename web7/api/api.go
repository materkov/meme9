package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	mongodriver "go.mongodb.org/mongo-driver/mongo"

	"github.com/materkov/meme9/web7/adapters/mongo"
)

const apiHost = ""
const staticHost = "/static"

func indexHTML() string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>meme9</title>
  <link rel="stylesheet" href="%s/index.css">
</head>
<body>
  <script>
    window.API_BASE_URL = "%s";
  </script>
  <div id="root"></div>
  <script src="%s/index.js"></script>
</body>
</html>`, staticHost, apiHost, staticHost)
}

type API struct {
	mongo *mongo.Adapter
}

func NewAPI(mongo *mongo.Adapter) *API {
	return &API{mongo: mongo}
}

type Post struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAd"`
}

func mapMongoPostToAPIPost(mongoPost mongo.Post, username string) Post {
	return Post{
		ID:        mongoPost.ID,
		Text:      mongoPost.Text,
		UserID:    mongoPost.UserID,
		Username:  username,
		CreatedAt: mongoPost.CreatedAt.Format(time.RFC3339),
	}
}

func (a *API) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func (a *API) feedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	mongoPosts, err := a.mongo.GetAllPosts(r.Context())
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Collect unique user IDs
	userIDSet := make(map[string]bool)
	for _, post := range mongoPosts {
		if post.UserID != "" {
			userIDSet[post.UserID] = true
		}
	}

	// Convert set to slice
	userIDs := make([]string, 0, len(userIDSet))
	for userID := range userIDSet {
		userIDs = append(userIDs, userID)
	}

	// Fetch all users in a single batch query
	users, err := a.mongo.GetUsersByIDs(r.Context(), userIDs)
	if err != nil {
		// Log error but continue with empty usernames
		log.Printf("Error fetching users: %v", err)
		users = make(map[string]*mongo.User)
	}

	// Build username map
	usernameMap := make(map[string]string)
	for userID, user := range users {
		usernameMap[userID] = user.Username
	}

	// Map posts to API posts with usernames
	apiPosts := make([]Post, len(mongoPosts))
	for i, mongoPost := range mongoPosts {
		username := usernameMap[mongoPost.UserID]
		apiPosts[i] = mapMongoPostToAPIPost(mongoPost, username)
	}

	json.NewEncoder(w).Encode(apiPosts)
}

type PublishReq struct {
	Text string `json:"text"`
}

type PublishResp struct {
	ID string `json:"id"`
}

func (a *API) verifyToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	// Support both "Bearer token" and just "token"
	tokenValue := strings.TrimPrefix(authHeader, "Bearer ")
	tokenValue = strings.TrimSpace(tokenValue)

	token, err := a.mongo.GetTokenByValue(r.Context(), tokenValue)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return "", fmt.Errorf("invalid token")
		}
		return "", fmt.Errorf("error verifying token: %w", err)
	}

	return token.UserID, nil
}

func (a *API) publishHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Verify authentication token
	userID, err := a.verifyToken(r)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	var publishReq PublishReq
	err = json.Unmarshal(body, &publishReq)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	post, err := a.mongo.AddPost(r.Context(), mongo.Post{
		Text:      publishReq.Text,
		UserID:    userID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create post"})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(PublishResp{ID: post.ID})
}

func (a *API) staticHandler(w http.ResponseWriter, r *http.Request) {
	// Strip /static prefix
	path := strings.TrimPrefix(r.URL.Path, "/static/")
	if path == "" {
		http.NotFound(w, r)
		return
	}

	// Build file path relative to web7 directory
	staticDir := filepath.Join("..", "front7", "dist")
	filePath := filepath.Join(staticDir, path)

	// Prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(staticDir)) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

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
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	var loginReq LoginReq
	err = json.Unmarshal(body, &loginReq)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "username and password required"})
		return
	}

	// Find user by username
	user, err := a.mongo.GetUserByUsername(r.Context(), loginReq.Username)
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
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

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(LoginResp{
		Token:    tokenValue,
		UserID:   user.ID,
		Username: user.Username,
	})
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

func (a *API) Serve() {
	http.HandleFunc("/feed", a.corsMiddleware(a.feedHandler))
	http.HandleFunc("/publish", a.corsMiddleware(a.publishHandler))
	http.HandleFunc("/login", a.corsMiddleware(a.loginHandler))
	http.HandleFunc("/register", a.corsMiddleware(a.registerHandler))
	http.HandleFunc("/static/", a.staticHandler)

	// Serve inline index.html from constant
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(indexHTML()))
	})

	log.Printf("Starting HTTP server at http://127.0.0.1:8080")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}
