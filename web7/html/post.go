package html

import (
	"fmt"
	"html"
	"time"
)

// PostPageData contains data for rendering a single post page
type PostPageData struct {
	PostID          string
	UserID          string
	Username        string
	Text            string
	CreatedAt       time.Time
	CurrentUsername string
}

// RenderPostPage renders the single post page HTML
func (r *Router) RenderPostPage(data PostPageData) string {
	formattedDate := data.CreatedAt.Format(time.RFC3339)
	escapedText := html.EscapeString(data.Text)

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
  <title>Post by %s</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
      background: #f5f5f5;
    }
    .container {
      background: transparent;
      border-radius: 8px;
      padding: 0;
    }
    .header {
      border-bottom: 1px solid #e0e0e0;
      padding-bottom: 20px;
      margin-bottom: 20px;
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
      display: inline-block;
      margin-bottom: 20px;
    }
    .backButton:hover {
      background: #e0e0e0;
      border-color: #ccc;
    }
    .post {
      background: #fff;
      border: 1px solid #e0e0e0;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
    .post .header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      border-bottom: 1px solid #e0e0e0;
      padding-bottom: 12px;
    }
    .post .header .username {
      font-size: 15px;
      font-weight: 600;
      color: #1976d2;
      text-decoration: none;
    }
    .post .header .username:hover {
      color: #1565c0;
      text-decoration: underline;
    }
    .post .header .date {
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
  </style>
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
</head>
<body>
  <a href="/feed" class="backButton">← Back to Feed</a>
  <div class="container">
    <div class="header">
      <h1>Post</h1>%s
    </div>
    <article class="post">
      <div class="header">
        <a href="/users/%s" class="username">%s</a>
        <time class="date">%s</time>
      </div>
      <p class="text">%s</p>
    </article>
  </div>
</body>
</html>`, data.Username, data.CurrentUsername, userInfoHTML, data.UserID, data.Username, formattedDate, escapedText)
}
