package api

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/materkov/meme9/web7/adapters/posts"
)

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

func (a *API) userPageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path /users/{id}
	userID := r.PathValue("id")
	if userID == "" {
		http.NotFound(w, r)
		return
	}

	// Fetch user info
	user, err := a.users.GetByID(r.Context(), userID)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		http.NotFound(w, r)
		return
	}

	username := user.Username
	if username == "" {
		username = "Unknown"
	}

	// Try to get current user and subscription status from cookie
	var isSubscribed bool
	cookie := r.Header.Get("Cookie")
	if cookie != "" {
		// Parse cookie to get auth_token
		authToken := ""
		cookies := parseCookies(cookie)
		if token, ok := cookies["auth_token"]; ok {
			authToken = token
		}

		if authToken != "" {
			currentUserIDFromToken, err := a.tokensService.VerifyToken(r.Context(), "Bearer "+authToken)
			if err == nil && currentUserIDFromToken != "" && currentUserIDFromToken != userID {
				// Check subscription status for SSR
				subscribed, err := a.subscriptions.IsSubscribed(r.Context(), currentUserIDFromToken, userID)
				if err == nil {
					isSubscribed = subscribed
				}
			}
		}
	}

	// Fetch posts for this user
	postsList, err := a.posts.GetByUserID(r.Context(), userID)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		postsList = []posts.Post{}
	}

	// Build posts HTML
	postsHTML := ""
	if len(postsList) == 0 {
		postsHTML = `<div class="empty">No posts yet</div>`
	} else {
		for _, post := range postsList {
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

	// Render HTML similar to frontend UserPostsPage
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>%s's Posts</title>
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
      align-items: center;
      gap: 1rem;
    }
    .backButton {
      padding: 0.5rem 1rem;
      background: #f5f5f5;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 0.9rem;
      cursor: pointer;
      text-decoration: none;
      color: #333;
      transition: all 0.2s;
    }
    .backButton:hover {
      background: #e0e0e0;
      border-color: #ccc;
    }
    .title {
      margin: 0;
      font-size: 32px;
      font-weight: 600;
      color: #333;
      flex: 1;
    }
    .subscribeSection {
      margin-left: auto;
    }
    .subscribeButton,
    .unsubscribeButton {
      padding: 0.5rem 1rem;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 0.9rem;
      cursor: pointer;
      transition: all 0.2s;
    }
    .subscribeButton {
      background: #1976d2;
      color: white;
      border-color: #1976d2;
    }
    .subscribeButton:hover:not(:disabled) {
      background: #1565c0;
      border-color: #1565c0;
    }
    .unsubscribeButton {
      background: #f5f5f5;
      color: #333;
    }
    .unsubscribeButton:hover:not(:disabled) {
      background: #e0e0e0;
      border-color: #ccc;
    }
    .subscribeButton:disabled,
    .unsubscribeButton:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
    .main {
      min-height: 400px;
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
    }
    .username {
      font-size: 15px;
      font-weight: 600;
      color: #1976d2;
      text-decoration: none;
    }
    .username:hover {
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
      <a href="/feed" class="backButton">← Back</a>
      <h1 class="title">%s's Posts</h1>
      <div class="subscribeSection" id="subscribeSection" style="display: none;">
        <button id="subscribeButton" class="subscribeButton" onclick="handleSubscribe()" style="display: none;">Subscribe</button>
        <button id="unsubscribeButton" class="unsubscribeButton" onclick="handleUnsubscribe()" style="display: none;">Unsubscribe</button>
      </div>
    </header>
    <main class="main">
      <div class="feed">
%s
      </div>
    </main>
  </div>

  <script>
    const userID = '%s';
    let isSubscribed = %t;
    let subscriptionLoading = false;

    // Load subscription status on page load
    (function() {
      const currentUserID = localStorage.getItem('auth_user_id');
      if (currentUserID && currentUserID !== userID) {
        // Show subscription section
        const subscribeSection = document.getElementById('subscribeSection');
        if (subscribeSection) {
          subscribeSection.style.display = 'block';
        }
        
        // Use server-provided subscription status
        updateSubscriptionUI();
      }
    })();

    function updateSubscriptionUI() {
      const subscribeBtn = document.getElementById('subscribeButton');
      const unsubscribeBtn = document.getElementById('unsubscribeButton');
      
      if (isSubscribed) {
        if (subscribeBtn) subscribeBtn.style.display = 'none';
        if (unsubscribeBtn) unsubscribeBtn.style.display = 'block';
      } else {
        if (subscribeBtn) subscribeBtn.style.display = 'block';
        if (unsubscribeBtn) unsubscribeBtn.style.display = 'none';
      }
    }

    async function handleSubscribe() {
      if (subscriptionLoading) return;
      
      const token = localStorage.getItem('auth_token');
      if (!token) {
        window.location.href = '/';
        return;
      }

      subscriptionLoading = true;
      const subscribeBtn = document.getElementById('subscribeButton');
      const unsubscribeBtn = document.getElementById('unsubscribeButton');
      
      if (subscribeBtn) {
        subscribeBtn.disabled = true;
        subscribeBtn.textContent = 'Subscribing...';
      }

      try {
        const response = await fetch('/api/subscribe', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
          },
          body: JSON.stringify({ user_id: userID })
        });

        if (response.ok) {
          isSubscribed = true;
          updateSubscriptionUI();
        } else {
          const error = await response.json();
          alert(error.error_details || 'Failed to subscribe');
        }
      } catch (err) {
        alert('Network error. Please try again.');
      } finally {
        subscriptionLoading = false;
        if (subscribeBtn) {
          subscribeBtn.disabled = false;
          subscribeBtn.textContent = 'Subscribe';
        }
      }
    }

    async function handleUnsubscribe() {
      if (subscriptionLoading) return;
      
      const token = localStorage.getItem('auth_token');
      if (!token) {
        window.location.href = '/';
        return;
      }

      subscriptionLoading = true;
      const subscribeBtn = document.getElementById('subscribeButton');
      const unsubscribeBtn = document.getElementById('unsubscribeButton');
      
      if (unsubscribeBtn) {
        unsubscribeBtn.disabled = true;
        unsubscribeBtn.textContent = 'Unsubscribing...';
      }

      try {
        const response = await fetch('/api/unsubscribe', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
          },
          body: JSON.stringify({ user_id: userID })
        });

        if (response.ok) {
          isSubscribed = false;
          updateSubscriptionUI();
        } else {
          const error = await response.json();
          alert(error.error_details || 'Failed to unsubscribe');
        }
      } catch (err) {
        alert('Network error. Please try again.');
      } finally {
        subscriptionLoading = false;
        if (unsubscribeBtn) {
          unsubscribeBtn.disabled = false;
          unsubscribeBtn.textContent = 'Unsubscribe';
        }
      }
    }
  </script>
</body>
</html>`, username, username, postsHTML, userID, isSubscribed)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}
