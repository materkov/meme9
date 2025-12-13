package render

import (
	"fmt"
	"html"
	"strings"
)

// FeedPageData contains data for rendering the feed page
type FeedPageData struct {
	FeedType              string
	Posts                 []*Post
	UsernameMap           map[string]string
	GlobalTabClass        string
	SubscriptionsTabClass string
	CurrentUsername       string
	IsAuthenticated       bool
}

// RenderFeedPage renders the feed page HTML
func RenderFeedPage(data FeedPageData) string {
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
			username := data.UsernameMap[post.UserId]
			if username == "" {
				username = "Unknown"
			}
			formattedDate := post.CreatedAt
			escapedText := html.EscapeString(post.Text)
			postsHTML += fmt.Sprintf(`
      <article class="post">
        <div class="header">
          <a href="/users/%s" class="username">%s</a>
          <time class="date">%s</time>
        </div>
        <p class="text"><a href="/posts/%s" class="post-link">%s</a></p>
      </article>`, post.UserId, username, formattedDate, post.Id, escapedText)
		}
	}

	// Build user info HTML
	userInfoHTML := `
        <div class="userInfo" id="userInfo" style="display: none;">
          <span class="username" id="currentUsername"></span>
          <button onclick="logout()" class="logout">Logout</button>
        </div>`

	// Build post form HTML (only show if authenticated)
	postFormHTML := ""
	if data.IsAuthenticated {
		postFormHTML = `<form class="postForm" id="postForm" onsubmit="handlePostSubmit(event)">
<textarea
  id="postText"
  placeholder="What's on your mind?"
  rows="4"
  maxlength="1000"
  required
></textarea>
<div class="footer">
  <div class="meta">
    <div id="postError" class="error" style="display: none;"></div>
    <div class="counter" id="postCounter">0 / 1000</div>
  </div>
  <button type="submit" class="button" id="postButton">Post</button>
</div>
</form>`
	}

	// Escape any % characters in all strings to prevent format specifier errors
	// This is needed because strings may contain % characters that would be interpreted as format specifiers
	userInfoHTMLStr := strings.ReplaceAll(userInfoHTML, "%", "%%")
	globalTabClassStr := strings.ReplaceAll(data.GlobalTabClass, "%", "%%")
	subscriptionsTabClassStr := strings.ReplaceAll(data.SubscriptionsTabClass, "%", "%%")
	postFormHTMLStr := strings.ReplaceAll(postFormHTML, "%", "%%")
	postsHTMLStr := strings.ReplaceAll(postsHTML, "%", "%%")
	currentUsernameStr := strings.ReplaceAll(data.CurrentUsername, "%", "%%")

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
      padding-bottom: 12px;
      position: relative;
      z-index: 1;
      overflow: visible;
      width: 100%%;
      box-sizing: border-box;
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
    .postForm {
      border: 1px solid #e0e0e0;
      border-radius: 8px;
      padding: 20px;
      margin-bottom: 20px;
      background: #fff;
    }
    .postForm textarea {
      width: 100%%;
      padding: 12px;
      border: 1px solid #e0e0e0;
      border-radius: 4px;
      font-size: 16px;
      font-family: inherit;
      resize: vertical;
      margin-bottom: 12px;
      box-sizing: border-box;
      min-height: 100px;
      background: #fff;
    }
    .postForm textarea:focus {
      outline: none;
      border-color: #666;
    }
    .postForm textarea:disabled {
      background-color: #f5f5f5;
      cursor: not-allowed;
    }
    .postForm .footer {
      display: flex;
      justify-content: space-between;
      align-items: flex-end;
      gap: 12px;
    }
    .postForm .meta {
      display: flex;
      flex-direction: column;
      gap: 4px;
      flex: 1;
    }
    .postForm .error {
      color: #dc3545;
      font-size: 14px;
      margin: 0;
    }
    .postForm .counter {
      font-size: 12px;
      color: #666;
      text-align: left;
    }
    .postForm .counter.error {
      color: #dc3545;
      font-weight: 500;
    }
    .postForm .button {
      padding: 10px 24px;
      background-color: #007bff;
      color: white;
      border: none;
      border-radius: 4px;
      font-size: 16px;
      font-weight: 500;
      cursor: pointer;
      transition: background-color 0.2s;
      white-space: nowrap;
    }
    .postForm .button:hover:not(:disabled) {
      background-color: #0056b3;
    }
    .postForm .button:disabled {
      background-color: #ccc;
      cursor: not-allowed;
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
%s
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

    // Post form handling
    const postForm = document.getElementById('postForm');
    const postText = document.getElementById('postText');
    const postCounter = document.getElementById('postCounter');
    const postError = document.getElementById('postError');
    const postButton = document.getElementById('postButton');

    if (postText && postCounter) {
      postText.addEventListener('input', function() {
        const length = this.value.length;
        postCounter.textContent = length + ' / 1000';
        if (length > 1000) {
          postCounter.classList.add('error');
        } else {
          postCounter.classList.remove('error');
        }
        if (postError) {
          postError.style.display = 'none';
        }
      });
    }

    async function handlePostSubmit(event) {
      event.preventDefault();
      
      if (!postText || !postButton) return;

      const text = postText.value.trim();
      if (!text || text.length > 1000) {
        return;
      }

      const token = getCookie('auth_token');
      if (!token) {
        window.location.href = '/';
        return;
      }

      postButton.disabled = true;
      postButton.textContent = 'Posting...';
      if (postError) {
        postError.style.display = 'none';
      }

      try {
        const response = await fetch('/twirp/meme.json_api.JsonAPI/Publish', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
          },
          body: JSON.stringify({ text: text })
        });

        if (response.ok) {
          // Clear form and reload page to show new post
          postText.value = '';
          postCounter.textContent = '0 / 1000';
          window.location.reload();
        } else {
          let error;
          try {
            error = await response.json();
          } catch (e) {
            error = { error: 'unknown_error', error_details: 'Failed to parse error response' };
          }
          const errorMsg = error.msg || error.error_details || 'Failed to create post';
          if (postError) {
            postError.textContent = errorMsg;
            postError.style.display = 'block';
          } else {
            alert(errorMsg);
          }
        }
      } catch (err) {
        const errorMsg = 'Network error. Please try again.';
        if (postError) {
          postError.textContent = errorMsg;
          postError.style.display = 'block';
        } else {
          alert(errorMsg);
        }
      } finally {
        postButton.disabled = false;
        postButton.textContent = 'Post';
      }
    }
  </script>
</body>
</html>`, userInfoHTMLStr, globalTabClassStr, subscriptionsTabClassStr, postFormHTMLStr, postsHTMLStr, currentUsernameStr)
}
