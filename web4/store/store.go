package store

type Post struct {
	ID int

	Text   string
	UserID int
}

type User struct {
	ID int

	Name      string
	LastPosts []int
}
