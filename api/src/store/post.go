package store

type Post struct {
	ID     int
	UserID int
	Date   int
	Text   string

	Link *PostLink

	IsDeleted bool
	PollID    int

	PhotoID int
}

type PostLink struct {
	Title       string
	Description string
	ImageURL    string
	URL         string
	FinalURL    string
}
