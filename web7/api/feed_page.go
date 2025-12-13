package api

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/materkov/meme9/web7/adapters/posts"
	"github.com/materkov/meme9/web7/adapters/users"
)

func (a *API) feedPageHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle GET requests (POST goes to API endpoint)
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Get feed type from query parameter (default to "global")
	feedType := r.URL.Query().Get("type")
	if feedType != "subscriptions" {
		feedType = "global"
	}

	// Try to get current user from cookie (optional - for subscriptions feed)
	var currentUserID string
	cookie := r.Header.Get("Cookie")
	if cookie != "" {
		cookies := parseCookies(cookie)
		if token, ok := cookies["auth_token"]; ok && token != "" {
			userID, err := a.tokensService.VerifyToken(r.Context(), "Bearer "+token)
			if err == nil {
				currentUserID = userID
			}
		}
	}

	// Fetch posts based on feed type
	var postsList []posts.Post
	var err error

	if feedType == "subscriptions" {
		if currentUserID == "" {
			// Redirect to login or show error
			http.Redirect(w, r, "/?error=Authentication required for subscriptions feed", http.StatusFound)
			return
		}

		// Get subscriptions for the current user
		followingIDs, err := a.subscriptions.GetFollowing(r.Context(), currentUserID)
		if err != nil {
			log.Printf("Error fetching subscriptions: %v", err)
			followingIDs = []string{}
		}

		// Include own posts and posts from subscribed users
		subscribedUserIDs := append(followingIDs, currentUserID)
		postsList, err = a.posts.GetByUserIDs(r.Context(), subscribedUserIDs)
		if err != nil {
			log.Printf("Error fetching subscription posts: %v", err)
			postsList = []posts.Post{}
		}
	} else {
		// Global feed - show all posts
		postsList, err = a.posts.GetAll(r.Context())
		if err != nil {
			log.Printf("Error fetching posts: %v", err)
			postsList = []posts.Post{}
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
	usersMap, err := a.users.GetByIDs(r.Context(), userIDs)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		usersMap = make(map[string]*users.User)
	}

	// Build username map
	usernameMap := make(map[string]string)
	for userID, user := range usersMap {
		usernameMap[userID] = user.Username
	}

	// Build posts HTML
	postsHTML := ""
	if len(postsList) == 0 {
		emptyMsg := "No posts yet"
		if feedType == "subscriptions" {
			emptyMsg = "No posts from your subscriptions yet"
		}
		postsHTML = fmt.Sprintf(`<div class="empty">%s</div>`, emptyMsg)
	} else {
		for _, post := range postsList {
			username := usernameMap[post.UserID]
			if username == "" {
				username = "Unknown"
			}
			formattedDate := post.CreatedAt.Format(time.RFC3339)
			escapedText := html.EscapeString(post.Text)
			postsHTML += fmt.Sprintf(`
      <article class="post">
        <div class="header">
          <a href="/users/%s" class="username">%s</a>
          <time class="date">%s</time>
        </div>
        <p class="text"><a href="/posts/%s" class="post-link">%s</a></p>
      </article>`, post.UserID, username, formattedDate, post.ID, escapedText)
		}
	}

	// Build user info HTML - will be populated by JavaScript from localStorage
	userInfoHTML := `
        <div class="userInfo" id="userInfo" style="display: none;">
          <span class="username" id="currentUsername"></span>
          <button onclick="logout()" class="logout">Logout</button>
        </div>`

	// Build feed tabs HTML
	globalTabClass := ""
	subscriptionsTabClass := ""
	if feedType == "global" {
		globalTabClass = "active"
	} else {
		subscriptionsTabClass = "active"
	}

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Posts Feed - meme9</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
      background: #f5f5f5;
    }
    .container {
      background: #fff;
      border-radius: 8px;
      padding: 20px;
    }
    .header {
      border-bottom: 1px solid #e0e0e0;
      padding-bottom: 20px;
      margin-bottom: 30px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .header h1 {
      margin: 0;
      font-size: 32px;
      font-weight: 600;
      color: #333;
    }
    .userInfo {
      display: flex;
      align-items: center;
      gap: 1rem;
    }
    .userInfo .username {
      color: #666;
      font-size: 0.9rem;
    }
    .logout {
      padding: 0.5rem 1rem;
      background: #f5f5f5;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 0.9rem;
      cursor: pointer;
      transition: all 0.2s;
    }
    .logout:hover {
      background: #e0e0e0;
      border-color: #ccc;
    }
    .main {
      min-height: 400px;
    }
    .feedTabs {
      display: flex;
      gap: 0.5rem;
      margin-bottom: 20px;
      border-bottom: 1px solid #e0e0e0;
    }
    .tab {
      padding: 0.75rem 1.5rem;
      background: transparent;
      border: none;
      border-bottom: 2px solid transparent;
      font-size: 1rem;
      cursor: pointer;
      color: #666;
      transition: all 0.2s;
      text-decoration: none;
      display: inline-block;
    }
    .tab:hover {
      color: #333;
      background: #f5f5f5;
    }
    .tab.active {
      color: #1976d2;
      border-bottom-color: #1976d2;
      font-weight: 600;
    }
    .empty {
      text-align: center;
      padding: 60px 20px;
      color: #666;
      font-size: 18px;
    }
    .feed {
      display: flex;
      flex-direction: column;
      gap: 20px;
    }
    .post {
      background: #fff;
      border: 1px solid #e0e0e0;
      border-radius: 8px;
      padding: 20px;
      transition: box-shadow 0.2s;
    }
    .post:hover {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
    .post .header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      border-bottom: 1px solid #e0e0e0;
      padding-bottom: 12px;
      margin-bottom: 12px;
      border: none;
      padding: 0;
    }
    .post .username {
      font-size: 15px;
      font-weight: 600;
      color: #1976d2;
      text-decoration: none;
    }
    .post .username:hover {
      color: #1565c0;
      text-decoration: underline;
    }
    .date {
      font-size: 14px;
      color: #999;
      font-weight: 400;
    }
    .text {
      margin: 0;
      font-size: 16px;
      line-height: 1.5;
      color: #333;
      white-space: pre-wrap;
      word-wrap: break-word;
    }
    .post-link {
      color: #333;
      text-decoration: none;
    }
    .post-link:hover {
      color: #1976d2;
    }
  </style>
</head>
<body>
  <div class="container">
    <header class="header">
      <h1>Posts Feed</h1>
      %s
    </header>
    <main class="main">
      <div class="feedTabs">
        <a href="/feed?type=global" class="tab %s">Global Feed</a>
        <a href="/feed?type=subscriptions" class="tab %s">Subscriptions</a>
      </div>
      <div class="feed">
%s
      </div>
    </main>
  </div>

  <script>
    // Load and display current user info
    (function() {
      const username = localStorage.getItem('auth_username');
      if (username) {
        const userInfoDiv = document.getElementById('userInfo');
        const usernameSpan = document.getElementById('currentUsername');
        if (userInfoDiv && usernameSpan) {
          usernameSpan.textContent = username;
          userInfoDiv.style.display = 'flex';
        }
      }
    })();

    function logout() {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('auth_user_id');
      localStorage.removeItem('auth_username');
      // Clear cookie
      document.cookie = 'auth_token=; path=/; max-age=0';
      window.location.href = '/';
    }
  </script>
</body>
</html>`, userInfoHTML, globalTabClass, subscriptionsTabClass, postsHTML)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}
