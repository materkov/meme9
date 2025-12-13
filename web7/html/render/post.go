package render

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
func RenderPostPage(data PostPageData) string {
	formattedDate := data.CreatedAt.Format(time.RFC3339)
	escapedText := html.EscapeString(data.Text)

	// Build page content with back button and post
	content := fmt.Sprintf(`<a href="/feed" class="backButton">← Back to Feed</a>
    <article class="post">
      <div class="header">
        <a href="/users/%s" class="username">%s</a>
        <time class="date">%s</time>
      </div>
      <p class="text">%s</p>
    </article>`, data.UserID, data.Username, formattedDate, escapedText)

	// Page-specific CSS
	extraCSS := `.container {
      background: transparent;
      padding: 0;
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
    }`

	// Page-specific JavaScript
	extraJS := GetUserInfoScript(data.CurrentUsername)

	// Use page container
	return RenderPageContainer(PageContainerData{
		Title:        fmt.Sprintf("Post by %s - meme9", data.Username),
		HeaderHTML:   "<h1>Post</h1>",
		UserInfoHTML: GetUserInfoHTML(),
		Content:      content,
		ExtraCSS:     extraCSS,
		ExtraJS:      extraJS,
	})
}
