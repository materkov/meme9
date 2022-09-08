package store

type Post struct {
	ID int

	Text   string
	Date   int
	UserID int
}

type User struct {
	ID   int
	VkID int

	Name string
}

type AuthToken struct {
	ID     int
	UserID int
	Token  string
	Date   int
}
