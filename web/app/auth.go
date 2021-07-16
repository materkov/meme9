package app

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/materkov/meme9/web/store"
)

var (
	ErrNotApplicable = fmt.Errorf("this auth not applicable")
	ErrAuthFailed    = fmt.Errorf("auth failed")
)

type Auth struct {
}

func (a *Auth) tryVkAuth(ctx context.Context, authUrl string) (int, error) {
	parsedUrl, err := url.Parse(authUrl)
	if err != nil || authUrl == "" {
		return 0, ErrNotApplicable
	}

	vkUserID, _ := strconv.Atoi(parsedUrl.Query().Get("vk_user_id"))
	if vkUserID == 0 {
		return 0, ErrNotApplicable
	}

	keys := make([]string, 0)
	for key := range parsedUrl.Query() {
		if strings.HasPrefix(key, "vk_") {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	for i, key := range keys {
		keys[i] = fmt.Sprintf("%s=%s", key, parsedUrl.Query().Get(key))
	}

	signString := strings.Join(keys, "&")

	mac := hmac.New(sha256.New, []byte(config.VKMiniAppSecret))
	_, _ = mac.Write([]byte(signString))
	computedSign := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	actualSign := parsedUrl.Query().Get("sign")
	if computedSign != actualSign {
		return 0, ErrAuthFailed
	}

	assocType := store.Assoc_VK_ID + strconv.Itoa(vkUserID)
	assocs, err := ObjectStore.AssocRange(ctx, 0, assocType, 1)
	if err != nil {
		return 0, fmt.Errorf("failed getting assoc: %w", err)
	}

	if len(assocs) > 0 {
		return assocs[0].VkID.ID2, nil
	}

	// New user
	userID, err := ObjectStore.GenerateNextID()
	if err != nil {
		return 0, fmt.Errorf("error generating object ID: %w", err)
	}

	user := store.User{
		ID:   userID,
		Name: fmt.Sprintf("VK User %d", vkUserID),
		VkID: vkUserID,
	}

	err = ObjectStore.ObjAdd(&store.StoredObject{ID: userID, User: &user})
	if err != nil {
		return 0, fmt.Errorf("error saving user: %w", err)
	}

	err = ObjectStore.AssocAdd(0, userID, assocType, &store.StoredAssoc{VkID: &store.VkID{
		ID1:  0,
		ID2:  userID,
		Type: assocType,
	}})
	if err != nil {
		return 0, fmt.Errorf("error saving assoc: %w", err)
	}

	return userID, nil
}

func (a *Auth) tryCookieAuth(r *http.Request) (*store.Token, error) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil || accessCookie.Value == "" {
		return nil, ErrNotApplicable
	}

	headerToken := r.Header.Get("x-csrf-token")
	validCSRFToken := GenerateCSRFToken(accessCookie.Value)

	if r.URL.Path == "/api" && headerToken != validCSRFToken {
		return nil, ErrAuthFailed
	}

	return a.tryTokenAuth(r.Context(), accessCookie.Value)
}

func (a *Auth) tryHeaderAuth(ctx context.Context, authHeader string) (*store.Token, error) {
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, ErrNotApplicable
	}

	return a.tryTokenAuth(ctx, authHeader)
}

func (a *Auth) tryTokenAuth(ctx context.Context, tokenStr string) (*store.Token, error) {
	tokenID := store.GetIdFromToken(tokenStr)
	if tokenID == 0 {
		return nil, ErrAuthFailed
	}

	obj, err := ObjectStore.ObjGet(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("error selecting token: %w", err)
	}

	if obj == nil || obj.Token == nil {
		return nil, ErrAuthFailed
	}
	if obj.Token.Token != tokenStr {
		return nil, ErrAuthFailed
	}

	return obj.Token, nil
}

func (a *Auth) TryAuth(r *http.Request) (*store.Token, int, error) {
	token, err := a.tryHeaderAuth(r.Context(), r.Header.Get("authorization"))
	if err == nil {
		return token, token.UserID, err
	}

	token, err = a.tryCookieAuth(r)
	if err == nil {
		return token, token.UserID, err
	}

	userID, err := a.tryVkAuth(r.Context(), r.URL.String())
	if err == nil {
		return nil, userID, err
	}

	userID, err = a.tryVkAuth(r.Context(), r.Header.Get("x-vk-auth"))
	if err == nil {
		return nil, userID, err
	}

	return nil, 0, fmt.Errorf("not authorized")
}
