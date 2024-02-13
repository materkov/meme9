package store2

type Store struct {
	Unique UniqueStore
	Likes  Likes
	Subs   Subscriptions
	Wall   Wall
	Votes  Votes

	Users       UserStore
	Posts       PostStore
	Polls       PollStore
	PollAnswers PollAnswerStore
	Tokens      TokenStore
	Bookmarks   BookmarkStore
	Files       *SqlFileStore
}

var GlobalStore *Store
