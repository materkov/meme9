package html

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/materkov/meme9/web7/api"
	json_api "github.com/materkov/meme9/web7/pb/github.com/materkov/meme9/api/json_api"
)

// Router handles HTML page routing
type Router struct {
	api *api.API
}

// NewRouter creates a new HTML router
func NewRouter(api *api.API) *Router {
	return &Router{api: api}
}

// parseCookies parses a cookie header string into a map
func parseCookies(cookieHeader string) map[string]string {
	cookies := make(map[string]string)
	pairs := strings.Split(cookieHeader, ";")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 {
			cookies[parts[0]] = parts[1]
		}
	}
	return cookies
}

// getAuthTokenFromRequest extracts the auth_token from the request cookies
func (r *Router) getAuthTokenFromRequest(req *http.Request) string {
	cookie := req.Header.Get("Cookie")
	if cookie == "" {
		return ""
	}
	cookies := parseCookies(cookie)
	if token, ok := cookies["auth_token"]; ok && token != "" {
		return token
	}
	return ""
}

// getCurrentUserID extracts and verifies the token from the request, returning the user ID if valid
func (r *Router) getCurrentUserID(req *http.Request) string {
	token := r.getAuthTokenFromRequest(req)
	if token == "" {
		return ""
	}
	// Use proto VerifyToken method
	verifyReq := &json_api.VerifyTokenRequest{
		Token: "Bearer " + token,
	}
	verifyResp, err := r.api.VerifyToken(req.Context(), verifyReq)
	if err != nil || verifyResp == nil {
		return ""
	}
	return verifyResp.UserId
}

// isAuthenticated checks if the request contains a valid authentication token
func (r *Router) isAuthenticated(req *http.Request) bool {
	return r.getCurrentUserID(req) != ""
}

// FeedPageHandler handles GET /feed requests
func (r *Router) FeedPageHandler(w http.ResponseWriter, req *http.Request) {
	// Only handle GET requests
	if req.Method != http.MethodGet {
		http.NotFound(w, req)
		return
	}

	// Get feed type from query parameter (default to "global")
	feedType := req.URL.Query().Get("type")
	if feedType != "subscriptions" {
		feedType = "global"
	}

	// Try to get current user from cookie (optional - for subscriptions feed)
	currentUserID := r.getCurrentUserID(req)

	// Fetch posts based on feed type using proto GetFeed
	feedReqType := "all"
	ctx := req.Context()
	if feedType == "subscriptions" {
		if currentUserID == "" {
			// Redirect to login or show error
			http.Redirect(w, req, "/?error=Authentication required for subscriptions feed", http.StatusFound)
			return
		}
		feedReqType = "subscriptions"
		// Set user ID in context for GetFeed to use
		ctx = context.WithValue(ctx, api.UserIDKey, currentUserID)
	}

	feedReq := &json_api.FeedRequest{
		Type: feedReqType,
	}
	feedResp, err := r.api.GetFeed(ctx, feedReq)
	var postsList []*Post
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		postsList = []*Post{}
	} else {
		// Convert FeedPostResponse to Post format
		postsList = make([]*Post, len(feedResp.Posts))
		for i, feedPost := range feedResp.Posts {
			postsList[i] = &Post{
				Id:        feedPost.Id,
				Text:      feedPost.Text,
				UserId:    feedPost.UserId,
				CreatedAt: feedPost.CreatedAt,
			}
		}
	}

	// Collect unique user IDs
	userIDSet := make(map[string]bool)
	for _, post := range postsList {
		if post.UserId != "" {
			userIDSet[post.UserId] = true
		}
	}

	// Convert set to slice
	userIDs := make([]string, 0, len(userIDSet))
	for userID := range userIDSet {
		userIDs = append(userIDs, userID)
	}

	// Fetch all users in a single batch query using proto GetUsersByIDs
	usersByIDsReq := &json_api.GetUsersByIDsRequest{
		UserIds: userIDs,
	}
	usersByIDsResp, err := r.api.GetUsersByIDs(req.Context(), usersByIDsReq)
	usersMap := make(map[string]*User)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
	} else {
		usersMap = usersByIDsResp.Users
	}

	// Build username map
	usernameMap := make(map[string]string)
	for userID, user := range usersMap {
		usernameMap[userID] = user.Username
	}

	// Build feed tabs HTML
	globalTabClass := ""
	subscriptionsTabClass := ""
	if feedType == "global" {
		globalTabClass = "active"
	} else {
		subscriptionsTabClass = "active"
	}

	// Get current user info for header using proto GetUser
	currentUsername := ""
	if currentUserID != "" {
		getUserReq := &json_api.GetUserRequest{
			UserId: currentUserID,
		}
		getUserResp, err := r.api.GetUser(req.Context(), getUserReq)
		if err == nil && getUserResp != nil {
			currentUsername = getUserResp.Username
		}
	}

	// Render HTML using html package
	htmlContent := r.RenderFeedPage(FeedPageData{
		FeedType:              feedType,
		Posts:                 postsList,
		UsernameMap:           usernameMap,
		GlobalTabClass:        globalTabClass,
		SubscriptionsTabClass: subscriptionsTabClass,
		CurrentUsername:       currentUsername,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}

// UserPageHandler handles GET /users/{id} requests
func (r *Router) UserPageHandler(w http.ResponseWriter, req *http.Request) {
	// Extract user ID from path /users/{id}
	userID := req.PathValue("id")
	if userID == "" {
		http.NotFound(w, req)
		return
	}

	ctx := req.Context()
	currentUserIDFromToken := r.getCurrentUserID(req)

	// Fetch user info, posts, and subscription status in parallel
	var (
		apiUser      *json_api.GetUserResponse
		apiPosts     []*json_api.GetPostResponse
		isSubscribed bool
		userErr      error
		postsErr     error
	)

	var wg sync.WaitGroup

	// Fetch user info using proto GetUser
	wg.Add(1)
	go func() {
		defer wg.Done()
		getUserReq := &json_api.GetUserRequest{
			UserId: userID,
		}
		apiUser, userErr = r.api.GetUser(ctx, getUserReq)
	}()

	// Fetch posts in parallel using proto GetUserPosts
	wg.Add(1)
	go func() {
		defer wg.Done()
		userPostsReq := &json_api.UserPostsRequest{
			UserId: userID,
		}
		userPostsResp, err := r.api.GetUserPosts(ctx, userPostsReq)
		if err != nil {
			postsErr = err
		} else {
			// Convert UserPostResponse to GetPostResponse format
			apiPosts = make([]*json_api.GetPostResponse, len(userPostsResp.Posts))
			for i, userPost := range userPostsResp.Posts {
				apiPosts[i] = &json_api.GetPostResponse{
					Id:        userPost.Id,
					Text:      userPost.Text,
					UserId:    userPost.UserId,
					CreatedAt: userPost.CreatedAt,
				}
			}
		}
	}()

	// Check subscription status in parallel using proto IsSubscribed (only if authenticated and different user)
	if currentUserIDFromToken != "" && currentUserIDFromToken != userID {
		wg.Add(1)
		go func() {
			defer wg.Done()
			isSubscribedReq := &json_api.IsSubscribedRequest{
				SubscriberId: currentUserIDFromToken,
				TargetUserId: userID,
			}
			isSubscribedResp, err := r.api.IsSubscribed(ctx, isSubscribedReq)
			if err == nil && isSubscribedResp != nil {
				isSubscribed = isSubscribedResp.Subscribed
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check if user exists
	if userErr != nil {
		log.Printf("Error fetching user: %v", userErr)
		http.NotFound(w, req)
		return
	}

	username := apiUser.Username
	if username == "" {
		username = "Unknown"
	}

	// Convert posts
	postsList := []*Post{}
	if postsErr != nil {
		log.Printf("Error fetching posts: %v", postsErr)
	} else {
		postsList = apiPosts
	}

	// Get current user info for header using proto GetUser
	currentUsername := ""
	if currentUserIDFromToken != "" {
		getUserReq := &json_api.GetUserRequest{
			UserId: currentUserIDFromToken,
		}
		getUserResp, err := r.api.GetUser(req.Context(), getUserReq)
		if err == nil && getUserResp != nil {
			currentUsername = getUserResp.Username
		}
	}

	// Determine if subscription section should be shown
	// Only show if authenticated and viewing someone else's profile
	showSubscribeSection := currentUserIDFromToken != "" && currentUserIDFromToken != userID

	// Render HTML using html package
	htmlContent := r.RenderUserPage(UserPageData{
		Username:             username,
		UserID:               userID,
		Posts:                postsList,
		IsSubscribed:         isSubscribed,
		CurrentUsername:      currentUsername,
		ShowSubscribeSection: showSubscribeSection,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}

// PostPageHandler handles GET /posts/{id} requests
func (r *Router) PostPageHandler(w http.ResponseWriter, req *http.Request) {
	// Extract post ID from path /posts/{id}
	postID := req.PathValue("id")
	if postID == "" {
		http.NotFound(w, req)
		return
	}

	// Fetch post from database using proto GetPost
	getPostReq := &json_api.GetPostRequest{
		PostId: postID,
	}
	apiPost, err := r.api.GetPost(req.Context(), getPostReq)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		http.NotFound(w, req)
		return
	}

	post := apiPost

	// Fetch user info using proto GetUser
	getUserReq := &json_api.GetUserRequest{
		UserId: post.UserId,
	}
	apiUser, err := r.api.GetUser(req.Context(), getUserReq)
	username := "Unknown"
	if err == nil && apiUser != nil {
		username = apiUser.Username
	}

	// Get current user info for header using proto GetUser
	currentUsername := ""
	currentUserID := r.getCurrentUserID(req)
	if currentUserID != "" {
		getCurrentUserReq := &json_api.GetUserRequest{
			UserId: currentUserID,
		}
		getCurrentUserResp, err := r.api.GetUser(req.Context(), getCurrentUserReq)
		if err == nil && getCurrentUserResp != nil {
			currentUsername = getCurrentUserResp.Username
		}
	}

	// Parse CreatedAt string to time.Time for template
	createdAt, _ := time.Parse(time.RFC3339, post.CreatedAt)

	// Render HTML using html package
	htmlContent := r.RenderPostPage(PostPageData{
		PostID:          post.Id,
		UserID:          post.UserId,
		Username:        username,
		Text:            post.Text,
		CreatedAt:       createdAt,
		CurrentUsername: currentUsername,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}

// AuthPageHandler handles GET / requests (auth page)
func (r *Router) AuthPageHandler(w http.ResponseWriter, req *http.Request) {
	// Get query parameter for tab (login or register)
	tab := req.URL.Query().Get("tab")
	if tab != "register" {
		tab = "login"
	}

	// Check for error messages from query params (for redirects after failed auth)
	errorMsg := req.URL.Query().Get("error")
	usernameError := req.URL.Query().Get("usernameError")
	credentialsError := req.URL.Query().Get("credentialsError")

	// Determine active tab classes
	loginTabClass := ""
	registerTabClass := ""
	if tab == "login" {
		loginTabClass = "active"
	} else {
		registerTabClass = "active"
	}

	// Submit button text
	submitText := "Login"
	if tab == "register" {
		submitText = "Register"
	}

	// Password autocomplete
	passwordAutocomplete := "current-password"
	if tab == "register" {
		passwordAutocomplete = "new-password"
	}

	// Determine input classes
	usernameInputClass := ""
	if usernameError != "" {
		usernameInputClass = `class="inputError"`
	}

	passwordInputClass := ""
	if credentialsError != "" {
		passwordInputClass = `class="inputError"`
	}

	// Render HTML using html package
	htmlContent := r.RenderAuthPage(AuthPageData{
		Tab:                  tab,
		UsernameError:        usernameError,
		CredentialsError:     credentialsError,
		Error:                errorMsg,
		LoginTabClass:        loginTabClass,
		RegisterTabClass:     registerTabClass,
		PasswordAutocomplete: passwordAutocomplete,
		SubmitText:           submitText,
		UsernameInputClass:   usernameInputClass,
		PasswordInputClass:   passwordInputClass,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}

// IndexHandler handles GET / requests (root route)
func (r *Router) IndexHandler(w http.ResponseWriter, req *http.Request) {
	// Check if user is authenticated - if yes, redirect to feed, otherwise show auth page
	if r.isAuthenticated(req) {
		// Redirect authenticated users to feed, preserving query parameters
		redirectURL := "/feed"
		if req.URL.RawQuery != "" {
			redirectURL += "?" + req.URL.RawQuery
		}
		http.Redirect(w, req, redirectURL, http.StatusFound)
		return
	}

	// Serve auth page for unauthenticated users
	r.AuthPageHandler(w, req)
}
