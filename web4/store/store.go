package store

type Store struct {
	Posts []Post
	Users []User
}

type Post struct {
	ID int

	Text   string
	UserID int
}

func (s *Store) GetPosts(ids []int) []Post {
	var result []Post
	for _, id := range ids {
		for _, post := range s.Posts {
			if post.ID == id {
				result = append(result, post)
			}
		}
	}

	return result
}

type User struct {
	ID int

	Name      string
	LastPosts []int
}

func (s *Store) GetUsers(ids []int) []User {
	var result []User
	for _, id := range ids {
		for _, user := range s.Users {
			if user.ID == id {
				result = append(result, user)
			}
		}
	}

	return result
}

var DefaultStore = Store{
	Posts: []Post{
		{
			ID:     100,
			Text:   "Post 100",
			UserID: 50,
		},
		{
			ID:     101,
			Text:   "Post 101",
			UserID: 51,
		},
		{
			ID:     102,
			Text:   "Post 102",
			UserID: 50,
		},
		{
			ID:     103,
			Text:   "Post 103",
			UserID: 51,
		},
	},
	Users: []User{
		{
			ID:        50,
			Name:      "User number patdisat",
			LastPosts: []int{102, 100},
		},
		{
			ID:        51,
			Name:      "User number patdisat adin",
			LastPosts: []int{103, 101},
		},
	},
}
