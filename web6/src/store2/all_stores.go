package store2

type Store struct {
	Unique     UniqueStore
	Nodes      NodeStore
	TypedNodes *TypedNodes
}

var GlobalStore *Store
