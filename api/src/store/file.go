package store

type File struct {
	ID     int
	UserID int

	Hash        string
	PhotoWidth  int
	PhotoHeight int
	Size        int
}
