package store

type Poll struct {
	ID        int
	UserID    int
	Question  string
	AnswerIds []int
}

func (p *Poll) GetAnswersMap() map[int]bool {
	result := map[int]bool{}
	for _, answerID := range p.AnswerIds {
		result[answerID] = true
	}

	return result
}

type PollAnswer struct {
	ID     int
	Answer string
}
