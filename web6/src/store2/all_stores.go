package store2

type Store struct {
	Unique     UniqueStore
	Nodes      NodeStore
	TypedNodes *TypedNodes
	Likes      Likes
	Subs       Subscriptions
	Wall       Wall
	Votes      Votes
}

var GlobalStore *Store
