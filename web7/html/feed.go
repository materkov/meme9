package html

import (
	"fmt"
	"html"
	"time"
)

// FeedPageData contains data for rendering the feed page
type FeedPageData struct {
	FeedType              string
	Posts                 []Post
	UsernameMap           map[string]string
	GlobalTabClass        string
	SubscriptionsTabClass string
	CurrentUsername       string
}

// RenderFeedPage renders the feed page HTML
func (r *Router) RenderFeedPage(data FeedPageData) string {
	// Build posts HTML
	postsHTML := ""
	if len(data.Posts) == 0 {
		emptyMsg := "No posts yet"
		if data.FeedType == "subscriptions" {
			emptyMsg = "No posts from your subscriptions yet"
		}
		postsHTML = fmt.Sprintf(`<div class="empty">%s</div>`, emptyMsg)
	} else {
		for _, post := range data.Posts {
			username := data.UsernameMap[post.UserID]
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

	// Build user info HTML
	userInfoHTML := `
        <div class="userInfo" id="userInfo" style="display: none;">
          <span class="username" id="currentUsername"></span>
          <button onclick="logout()" class="logout">Logout</button>
        </div>`

	return fmt.Sprintf(`<!DOCTYPE html>
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
      const username = '%s';
      if (username) {
        const userInfoDiv = document.getElementById('userInfo');
        const usernameSpan = document.getElementById('currentUsername');
        if (userInfoDiv && usernameSpan) {
          usernameSpan.textContent = username;
          userInfoDiv.style.display = 'flex';
        }
      }
    })();

    function getCookie(name) {
      const value = '; ' + document.cookie;
      const parts = value.split('; ' + name + '=');
      if (parts.length === 2) return parts.pop().split(';').shift();
      return null;
    }

    function logout() {
      // Clear cookie
      document.cookie = 'auth_token=; path=/; max-age=0';
      window.location.href = '/';
    }
  </script>
</body>
</html>`, userInfoHTML, data.GlobalTabClass, data.SubscriptionsTabClass, postsHTML, data.CurrentUsername)
}
