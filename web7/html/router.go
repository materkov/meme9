package html

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/materkov/meme9/web7/api"
)

// Router handles HTML page routing
type Router struct {
	api *api.API
}

// NewRouter creates a new HTML router
func NewRouter(api *api.API) *Router {
	return &Router{api: api}
}

// convertAPIPostToPost converts API PostData to html Post
func convertAPIPostToPost(apiPost api.PostData) Post {
	return Post{
		ID:        apiPost.ID,
		Text:      apiPost.Text,
		UserID:    apiPost.UserID,
		CreatedAt: apiPost.CreatedAt,
	}
}

// convertAPIPostsToPosts converts API PostData slice to html Post slice
func convertAPIPostsToPosts(apiPosts []api.PostData) []Post {
	result := make([]Post, len(apiPosts))
	for i, apiPost := range apiPosts {
		result[i] = convertAPIPostToPost(apiPost)
	}
	return result
}

// convertAPIUserToUser converts API UserData to html User
func convertAPIUserToUser(apiUser *api.UserData) *User {
	if apiUser == nil {
		return nil
	}
	return &User{
		ID:       apiUser.ID,
		Username: apiUser.Username,
	}
}

// convertAPIUsersToUsers converts API UserData map to html User map
func convertAPIUsersToUsers(apiUsers map[string]*api.UserData) map[string]*User {
	result := make(map[string]*User)
	for userID, apiUser := range apiUsers {
		result[userID] = convertAPIUserToUser(apiUser)
	}
	return result
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
	userID, err := r.api.VerifyToken(req.Context(), "Bearer "+token)
	if err != nil {
		return ""
	}
	return userID
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

	// Fetch posts based on feed type
	var postsList []Post
	var err error

	if feedType == "subscriptions" {
		if currentUserID == "" {
			// Redirect to login or show error
			http.Redirect(w, req, "/?error=Authentication required for subscriptions feed", http.StatusFound)
			return
		}

		// Get subscriptions for the current user
		followingIDs, err := r.api.GetFollowing(req.Context(), currentUserID)
		if err != nil {
			log.Printf("Error fetching subscriptions: %v", err)
			followingIDs = []string{}
		}

		// Include own posts and posts from subscribed users
		subscribedUserIDs := append(followingIDs, currentUserID)
		apiPosts, err := r.api.GetPostsByUserIDsHTML(req.Context(), subscribedUserIDs)
		if err != nil {
			log.Printf("Error fetching subscription posts: %v", err)
			postsList = []Post{}
		} else {
			postsList = convertAPIPostsToPosts(apiPosts)
		}
	} else {
		// Global feed - show all posts
		apiPosts, err := r.api.GetAllPostsHTML(req.Context())
		if err != nil {
			log.Printf("Error fetching posts: %v", err)
			postsList = []Post{}
		} else {
			postsList = convertAPIPostsToPosts(apiPosts)
		}
	}

	// Collect unique user IDs
	userIDSet := make(map[string]bool)
	for _, post := range postsList {
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
	apiUsersMap, err := r.api.GetUsersByIDsHTML(req.Context(), userIDs)
	usersMap := make(map[string]*User)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
	} else {
		usersMap = convertAPIUsersToUsers(apiUsersMap)
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

	// Get current user info for header
	currentUsername := ""
	if currentUserID != "" {
		currentUser, err := r.api.GetUserByIDHTML(req.Context(), currentUserID)
		if err == nil && currentUser != nil {
			currentUsername = currentUser.Username
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
		apiUser      *api.UserData
		apiPosts     []api.PostData
		isSubscribed bool
		userErr      error
		postsErr     error
	)

	var wg sync.WaitGroup

	// Fetch user info
	wg.Add(1)
	go func() {
		defer wg.Done()
		apiUser, userErr = r.api.GetUserByIDHTML(ctx, userID)
	}()

	// Fetch posts in parallel
	wg.Add(1)
	go func() {
		defer wg.Done()
		apiPosts, postsErr = r.api.GetPostsByUserIDHTML(ctx, userID)
	}()

	// Check subscription status in parallel (only if authenticated and different user)
	if currentUserIDFromToken != "" && currentUserIDFromToken != userID {
		wg.Add(1)
		go func() {
			defer wg.Done()
			subscribed, err := r.api.IsSubscribed(ctx, currentUserIDFromToken, userID)
			if err == nil {
				isSubscribed = subscribed
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
	postsList := []Post{}
	if postsErr != nil {
		log.Printf("Error fetching posts: %v", postsErr)
	} else {
		postsList = convertAPIPostsToPosts(apiPosts)
	}

	// Get current user info for header
	currentUsername := ""
	if currentUserIDFromToken != "" {
		currentUser, err := r.api.GetUserByIDHTML(req.Context(), currentUserIDFromToken)
		if err == nil && currentUser != nil {
			currentUsername = currentUser.Username
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

	// Fetch post from database
	apiPost, err := r.api.GetPostByIDHTML(req.Context(), postID)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		http.NotFound(w, req)
		return
	}

	post := convertAPIPostToPost(*apiPost)

	// Fetch user info
	apiUser, err := r.api.GetUserByIDHTML(req.Context(), post.UserID)
	username := "Unknown"
	if err == nil && apiUser != nil {
		username = apiUser.Username
	}

	// Get current user info for header
	currentUsername := ""
	currentUserID := r.getCurrentUserID(req)
	if currentUserID != "" {
		currentUser, err := r.api.GetUserByIDHTML(req.Context(), currentUserID)
		if err == nil && currentUser != nil {
			currentUsername = currentUser.Username
		}
	}

	// Render HTML using html package
	htmlContent := r.RenderPostPage(PostPageData{
		PostID:          post.ID,
		UserID:          post.UserID,
		Username:        username,
		Text:            post.Text,
		CreatedAt:       post.CreatedAt,
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
