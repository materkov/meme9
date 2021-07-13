package store

const (
	Assoc_Liked = 1
)

type StoredAssoc struct {
	Liked   *Liked
	LikedBy *LikedBy
}

type Liked struct {
}

type LikedBy struct {
}
