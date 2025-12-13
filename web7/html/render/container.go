package render

import (
	"fmt"
	"strings"
)

// PageContainerData contains data for the page container
type PageContainerData struct {
	Title        string
	HeaderHTML   string // Custom header HTML (if empty, uses Title in h1)
	UserInfoHTML string
	Content      string
	ExtraCSS     string
	ExtraJS      string
}

// RenderPageContainer renders a page with common structure, CSS, and JavaScript
func RenderPageContainer(data PageContainerData) string {
	// Escape % characters to prevent format specifier errors
	titleStr := strings.ReplaceAll(data.Title, "%", "%%")
	headerHTMLStr := strings.ReplaceAll(data.HeaderHTML, "%", "%%")
	userInfoHTMLStr := strings.ReplaceAll(data.UserInfoHTML, "%", "%%")
	contentStr := strings.ReplaceAll(data.Content, "%", "%%")
	extraCSSStr := strings.ReplaceAll(data.ExtraCSS, "%", "%%")
	extraJSStr := strings.ReplaceAll(data.ExtraJS, "%", "%%")

	// Use custom header HTML if provided, otherwise use title in h1
	headerContent := headerHTMLStr
	if headerContent == "" {
		headerContent = fmt.Sprintf(`<h1>%s</h1>`, titleStr)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>%s</title>
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
    %s
  </style>
</head>
<body>
  <div class="container">
    <header class="header">
      %s
      %s
    </header>
    <main class="main">
      %s
    </main>
  </div>

  <script>
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

    %s
  </script>
</body>
</html>`, titleStr, extraCSSStr, headerContent, userInfoHTMLStr, contentStr, extraJSStr)
}

// GetUserInfoHTML returns the standard user info HTML snippet
func GetUserInfoHTML() string {
	return `
        <div class="userInfo" id="userInfo" style="display: none;">
          <span class="username" id="currentUsername"></span>
          <button onclick="logout()" class="logout">Logout</button>
        </div>`
}

// GetUserInfoScript returns JavaScript to set username in user info
func GetUserInfoScript(username string) string {
	if username == "" {
		return ""
	}
	return fmt.Sprintf(`
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
    })();`, strings.ReplaceAll(username, "'", "\\'"))
}
