package store

const (
	Assoc_Liked     = 1
	Assoc_Commended = 2
	Assoc_Following = 3
)

type StoredAssoc struct {
	Liked     *Liked
	Commented *Commented
	Following *Following
}

type Liked struct {
	ID1  int
	ID2  int
	Type int
}

type Commented struct {
	ID1  int
	ID2  int
	Type int
}

type Following struct {
	ID1  int
	ID2  int
	Type int
}
