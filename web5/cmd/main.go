package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/materkov/meme9/web5/api"
	"github.com/materkov/meme9/web5/imgproxy"
	"github.com/materkov/meme9/web5/pkg/auth"
	"github.com/materkov/meme9/web5/pkg/files"
	"github.com/materkov/meme9/web5/pkg/telegram"
	"github.com/materkov/meme9/web5/pkg/users"
	"github.com/materkov/meme9/web5/store"
	"github.com/materkov/meme9/web5/upload"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ApiError string

func (e ApiError) Error() string {
	return string(e)
}

func write(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)

		var apiErr ApiError
		if errors.As(err, &apiErr) {
			_, _ = fmt.Fprint(w, err.Error())
		} else if err != nil {
			log.Printf("[ERROR] Internal error: %s", err)
			_, _ = fmt.Fprint(w, "internal error")
		}
	} else {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func handleSetOnline(w http.ResponseWriter, r *http.Request) {
	viewer := r.Context().Value(ViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		write(w, nil, ApiError("not authorized"))
		return
	}

	go func() {
		_, err := store.RedisClient.Set(context.Background(), fmt.Sprintf("online:%d", viewer.UserID), time.Now().Unix(), time.Minute*3).Result()
		if err != nil {
			log.Printf("Err: %s", err)
		}
	}()

	write(w, nil, nil)
}

func handleUserEdit(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.FormValue("id"))

	user := store.User{}
	err := store.NodeGet(userID, &user)
	if err != nil {
		write(w, nil, err)
		return
	} else if user.ID == 0 {
		write(w, nil, ApiError("user not found"))
		return
	}

	viewer := r.Context().Value(ViewerKey).(*Viewer)
	if viewer.UserID != user.ID {
		write(w, nil, ApiError("no access to edit this user"))
		return
	}

	name := r.FormValue("name")
	if name == "" {
		write(w, nil, ApiError("name is empty"))
		return
	} else if len(name) > 100 {
		write(w, nil, ApiError("name is too long"))
		return
	}

	user.Name = name

	err = store.NodeSave(user.ID, user)
	if err != nil {
		write(w, nil, err)
		return
	}

	write(w, nil, nil)
}

func handleVkCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirectUri")

	vkID, vkAccessToken, err := auth.ExchangeCode(code, redirectURI)
	if err != nil {
		write(w, nil, err)
		return
	}

	userID, err := users.GetOrCreateByVKID(vkID)
	if err != nil {
		write(w, nil, err)
		return
	}

	user := &store.User{}
	err = store.NodeGet(userID, user)
	if err != nil {
		write(w, nil, err)
		return
	}

	user.VkAccessToken = vkAccessToken
	err = store.NodeSave(user.ID, user)
	if err != nil {
		log.Printf("error saving user")
	}

	_, _ = store.RedisClient.RPush(context.Background(), "queue", user.ID).Result()

	authToken, err := auth.CreateToken(userID)
	if err != nil {
		write(w, nil, err)
		return
	}

	err = telegram.SendNotify(fmt.Sprintf("meme new login: https://vk.com/id%d", user.VkID))
	if err != nil {
		log.Printf("Error sending telegram notify: %s", err)
	}

	resp := []interface{}{
		authToken,
		userID,
	}
	write(w, resp, nil)
}

func handleEmailRegister(w http.ResponseWriter, r *http.Request) {
	email, password := r.FormValue("email"), r.FormValue("password")
	validateErr := auth.ValidateCredentials(email, password)
	if validateErr != "" {
		write(w, nil, ApiError(validateErr))
		return
	}

	userID, err := auth.Register(email, password)
	if err != nil {
		write(w, nil, err)
		return
	}

	token, err := auth.CreateToken(userID)
	if err != nil {
		write(w, nil, err)
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	write(w, resp, nil)
}

func handleAuthEmail(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.EmailAuth(r.FormValue("email"), r.FormValue("password"))
	if err == auth.ErrInvalidCredentials {
		write(w, nil, ApiError("invalid credentials"))
		return
	} else if err != nil {
		write(w, nil, err)
		return
	}

	token, err := auth.CreateToken(userID)
	if err != nil {
		write(w, nil, err)
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	write(w, resp, nil)
}

func handleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		write(w, nil, ApiError("invalid file"))
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		write(w, nil, err)
		return
	}

	hash := sha256.Sum256(fileBytes)
	hashHex := hex.EncodeToString(hash[:])

	err = files.SelectelUpload(fileBytes, hashHex)
	if err != nil {
		write(w, nil, err)
		return
	}

	viewer := r.Context().Value(ViewerKey).(*Viewer)
	if viewer.UserID == 0 {
		write(w, nil, ApiError("not authorized"))
		return
	}

	user := &store.User{}
	err = store.NodeGet(viewer.UserID, user)
	if err != nil {
		write(w, nil, err)
		return
	}

	user.AvatarSha = hashHex

	err = store.NodeSave(user.ID, user)
	if err != nil {
		write(w, nil, err)
		return
	}

	resp := struct {
		Avatar string `json:"avatar"`
	}{
		Avatar: files.GetURL(hashHex),
	}
	write(w, resp, nil)
}

type Viewer struct {
	UserID int

	Origin string
}

func (v *Viewer) GetUserIDStr() string {
	if v.UserID == 0 {
		return ""
	}
	return strconv.Itoa(v.UserID)
}

type contextKey int

const ViewerKey contextKey = iota

func wrapper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			return
		}

		authToken := r.FormValue("token")
		if authToken == "" {
			authToken = r.Header.Get("Authorization")
			authToken = strings.TrimPrefix(authToken, "Bearer ")
		}
		userID, _ := auth.CheckToken(authToken)

		viewer := &Viewer{
			UserID: userID,
			Origin: r.Header.Get("origin"),
		}
		r = r.WithContext(context.WithValue(r.Context(), ViewerKey, viewer))

		next(w, r)
	}
}

func main() {
	queue := ""
	flag.StringVar(&queue, "queue", "", "Queue listen to")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	store.RedisClient = redis.NewClient(&redis.Options{})

	configStr, err := store.RedisClient.Get(context.Background(), "config").Bytes()
	if err != nil {
		log.Fatalf("Failed reading config: %s", err)
	}

	err = json.Unmarshal(configStr, &store.DefaultConfig)
	if err != nil {
		log.Fatalf("Error parsing config JSON: %s", err)
	}

	if queue != "" {
		HandleWorker(queue)
		return
	}

	http.HandleFunc("/api", api.HandleAPI)
	http.HandleFunc("/api2/", api.HandleAPI2)
	http.HandleFunc("/upload", upload.HandleUpload)
	http.HandleFunc("/imgproxy", imgproxy.ServeHTTP)

	http.HandleFunc("/api/setOnline", wrapper(handleSetOnline))
	http.HandleFunc("/api/userEdit", wrapper(handleUserEdit))
	http.HandleFunc("/api/vkCallback", wrapper(handleVkCallback))
	http.HandleFunc("/api/emailRegister", wrapper(handleEmailRegister))
	http.HandleFunc("/api/emailLogin", wrapper(handleAuthEmail))
	http.HandleFunc("/api/uploadAvatar", wrapper(handleUploadAvatar))

	_ = http.ListenAndServe("127.0.0.1:8000", nil)
}

func HandleWorker(queue string) {
	for {
		result, err := store.RedisClient.BLPop(context.Background(), time.Second*5, queue).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return
		}

		log.Printf("Got queue task: %v", result)

		userID, _ := strconv.Atoi(result[1])
		err = users.RefreshFromVk(userID)
		if err != nil {
			log.Printf("Error doing queue: %s", err)
		}
	}
}
