package render

import (
	"fmt"
	"html"
)

// UserPageData contains data for rendering the user posts page
type UserPageData struct {
	Username             string
	UserID               string
	Posts                []*Post
	IsSubscribed         bool
	CurrentUsername      string
	ShowSubscribeSection bool
}

// RenderUserPage renders the user posts page HTML
func RenderUserPage(data UserPageData) string {
	// Determine display style for subscription section
	subscribeSectionDisplay := "none"
	if data.ShowSubscribeSection {
		subscribeSectionDisplay = "block"
	}

	// Build posts HTML
	postsHTML := ""
	if len(data.Posts) == 0 {
		postsHTML = `<div class="empty">No posts yet</div>`
	} else {
		for _, post := range data.Posts {
			formattedDate := post.CreatedAt
			escapedText := html.EscapeString(post.Text)
			postsHTML += fmt.Sprintf(`
      <article class="post">
        <div class="header">
          <a href="/users/%s" class="username">%s</a>
          <time class="date">%s</time>
        </div>
        <p class="text"><a href="/posts/%s" class="post-link">%s</a></p>
      </article>`, post.UserId, data.Username, formattedDate, post.Id, escapedText)
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
      justify-content: space-between;
      align-items: center;
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
      display: inline-block;
      margin-bottom: 20px;
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
  <a href="/feed" class="backButton">← Back to Feed</a>
  <div class="container">
    <header class="header">
      <h1 class="title">%s's Posts</h1>
      <div class="subscribeSection" id="subscribeSection" style="display: %s;">
        <button id="subscribeButton" class="subscribeButton" onclick="handleSubscribe()" style="display: none;">Subscribe</button>
        <button id="unsubscribeButton" class="unsubscribeButton" onclick="handleUnsubscribe()" style="display: none;">Unsubscribe</button>
      </div>%s
    </header>
    <main class="main">
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

    const userID = '%s';
    let isSubscribed = %t;
    let subscriptionLoading = false;

    // Load subscription status on page load
    (function() {
      // Subscription section visibility is controlled by server
      // Only update UI if section is visible
      const subscribeSection = document.getElementById('subscribeSection');
      if (subscribeSection && subscribeSection.style.display !== 'none') {
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
      
      const token = getCookie('auth_token');
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
        const response = await fetch('/twirp/meme.json_api.JsonAPI/Subscribe', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
          },
          body: JSON.stringify({ user_id: userID })
        });

        if (response.ok) {
          const result = await response.json();
          isSubscribed = result.subscribed || true;
          updateSubscriptionUI();
        } else {
          let error;
          try {
            error = await response.json();
          } catch (e) {
            error = { error: 'unknown_error', error_details: 'Failed to parse error response' };
          }
          alert(error.msg || error.error_details || 'Failed to subscribe');
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
      
      const token = getCookie('auth_token');
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
        const response = await fetch('/twirp/meme.json_api.JsonAPI/Unsubscribe', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
          },
          body: JSON.stringify({ user_id: userID })
        });

        if (response.ok) {
          const result = await response.json();
          isSubscribed = result.subscribed || false;
          updateSubscriptionUI();
        } else {
          let error;
          try {
            error = await response.json();
          } catch (e) {
            error = { error: 'unknown_error', error_details: 'Failed to parse error response' };
          }
          alert(error.msg || error.error_details || 'Failed to unsubscribe');
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
</html>`, data.Username, data.Username, subscribeSectionDisplay, userInfoHTML, postsHTML, data.CurrentUsername, data.UserID, data.IsSubscribed)
}
