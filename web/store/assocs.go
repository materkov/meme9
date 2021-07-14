package store

const (
	Assoc_Liked     = "Liked"
	Assoc_Commended = "Commented"
	Assoc_Following = "Following"
	Assoc_VK_ID     = "VkID:"
	AssocPosted     = "Posted"

	FakeIDVK = 1045
)

type StoredAssoc struct {
	Liked     *Liked
	Commented *Commented
	Following *Following
	VkID      *VkID
	Posted    *Posted
}

type Liked struct {
	ID1  int
	ID2  int
	Type string
}

type Commented struct {
	ID1  int
	ID2  int
	Type string
}

type Following struct {
	ID1  int
	ID2  int
	Type string
}

type VkID struct {
	ID1  int
	ID2  int
	Type string
}

type Posted struct {
	ID1  int
	ID2  int
	Type string
}
