package store

type StoredObject struct {
	ID     int
	APILog *APILog
	Photo  *Photo
}

type APILog struct {
	ID       int
	UserID   int
	Method   string
	Request  string
	Response string
}

type Photo struct {
	ID     int
	UserID int
	Path   string
}

type User struct {
	ID   int
	VkID int
}
