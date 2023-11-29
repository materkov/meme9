package store

type Poll struct {
	ID        int
	UserID    int
	Question  string
	AnswerIds []int
}

type PollAnswer struct {
	ID     int
	Answer string
}
