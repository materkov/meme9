package api

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

func (a *API) postPageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract post ID from path /posts/{id}
	postID := r.PathValue("id")
	if postID == "" {
		http.NotFound(w, r)
		return
	}

	// Fetch post from database
	post, err := a.posts.GetByID(r.Context(), postID)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		http.NotFound(w, r)
		return
	}

	// Fetch user info
	user, err := a.users.GetByID(r.Context(), post.UserID)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
	}

	username := "Unknown"
	if user != nil {
		username = user.Username
	}

	// Format date
	formattedDate := post.CreatedAt.Format(time.RFC3339)

	// Escape HTML in text to prevent XSS
	escapedText := html.EscapeString(post.Text)

	// Render HTML similar to frontend Post component
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
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
    .post {
      background: #fff;
      border: 1px solid #e0e0e0;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
    .header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      border-bottom: 1px solid #e0e0e0;
      padding-bottom: 12px;
    }
    .username {
      font-size: 15px;
      font-weight: 600;
      color: #1976d2;
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
    .back-link {
      display: inline-block;
      margin-bottom: 20px;
      color: #1976d2;
      text-decoration: none;
    }
    .back-link:hover {
      text-decoration: underline;
    }
  </style>
</head>
<body>
  <a href="/users/%s" class="back-link">← Back to %s's posts</a>
  <article class="post">
    <div class="header">
      <a href="/users/%s" class="username">%s</a>
      <time class="date">%s</time>
    </div>
    <p class="text">%s</p>
  </article>
</body>
</html>`, username, post.UserID, username, post.UserID, username, formattedDate, escapedText)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}
