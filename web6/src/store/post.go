package store

type Post struct {
	ID     int
	UserID int
	Date   int
	Text   string

	Link *PostLink
}

type PostLink struct {
	Title       string
	Description string
	ImageURL    string
	URL         string
	FinalURL    string
}
