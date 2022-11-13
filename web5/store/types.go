package store

type Post struct {
	ID     int
	Date   int
	Text   string
	UserID int

	IsDeleted bool
	PhotoID   int
}

type User struct {
	ID   int
	Name string
	Bio  string

	VkID          int
	VkAccessToken string
	VkPhoto200    string

	Email        string
	PasswordHash string

	AvatarSha string
}

type AuthToken struct {
	ID     int
	UserID int
	Token  string
	Date   int
}

type Config struct {
	VKAppID     int
	VKAppSecret string

	TelegramToken string

	SelectelAccountID    int
	SelectelUserName     string
	SelectelUserPassword string
}

type Photo struct {
	ID int

	Size   int
	Hash   string
	Width  int
	Height int
}
