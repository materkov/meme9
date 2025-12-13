package html

import "time"

// Post represents a post for HTML rendering
type Post struct {
	ID        string
	Text      string
	UserID    string
	CreatedAt time.Time
}

// User represents a user for HTML rendering
type User struct {
	ID       string
	Username string
}
