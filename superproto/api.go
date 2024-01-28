package api

import "context"

type FeedType int

const (
	FeedType_FEED     FeedType = 0
	FeedType_DISCOVER FeedType = 1
)

type PostLikeAction int

const (
	PostLikeAction_LIKE   PostLikeAction = 0
	PostLikeAction_UNLIKE PostLikeAction = 1
)

type SubscribeAction int

func (s SubscribeAction) MarshalJSON() ([]byte, error) {
	m := map[SubscribeAction][]byte{
		SubscribeAction_FOLLOW:   []byte("\"FOLLOW\""),
		SubscribeAction_UNFOLLOW: []byte("\"UNFOLLOW\""),
	}
	return m[s], nil
}

const (
	SubscribeAction_FOLLOW   SubscribeAction = 0
	SubscribeAction_UNFOLLOW SubscribeAction = 1
)

type Posts interface {
	Add(context.Context, *AddReq) (*Post, error)
	List(context.Context, *ListReq) (*PostsList, error)
	Delete(context.Context, *PostsDeleteReq) (*Void, error)
	Like(context.Context, *PostsLikeReq) (*Void, error)
}

type Void struct {
}

type AddReq struct {
	Text   string `json:"text,omitempty"`
	PollId string `json:"pollId,omitempty"`
}

type Post struct {
	Id         string    `json:"id,omitempty"`
	UserId     string    `json:"userId,omitempty"`
	Date       string    `json:"date,omitempty"`
	Text       string    `json:"text,omitempty"`
	User       *User     `json:"user,omitempty"`
	IsLiked    bool      `json:"isLiked,omitempty"`
	LikesCount int       `json:"likesCount,omitempty"`
	Link       *PostLink `json:"link,omitempty"`
	Poll       *Poll     `json:"poll,omitempty"`
}

type PostLink struct {
	Url         string `json:"url,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ImageUrl    string `json:"imageUrl,omitempty"`
	Domain      string `json:"domain,omitempty"`
}

type PostsList struct {
	Items     []*Post `json:"items,omitempty"`
	PageToken string  `json:"pageToken,omitempty"`
}

type ListReq struct {
	Type      FeedType `json:"type,omitempty"`
	ByUserId  string   `json:"byUserId,omitempty"`
	ById      string   `json:"byId,omitempty"`
	Count     int      `json:"count,omitempty"`
	PageToken string   `json:"pageToken,omitempty"`
}

type PostsDeleteReq struct {
	PostId string `json:"postId,omitempty"`
}

type PostsLikeReq struct {
	PostId string         `json:"postId,omitempty"`
	Action PostLikeAction `json:"action,omitempty"`
}

type Users interface {
	List(context.Context, *UsersListReq) (*UsersList, error)
	SetStatus(context.Context, *UsersSetStatus) (*Void, error)
	Follow(context.Context, *UsersFollowReq) (*Void, error)
}

type User struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
	IsFollowing bool   `json:"isFollowing,omitempty"`
}

type UsersList struct {
	Users []*User `json:"users,omitempty"`
}

type UsersListReq struct {
	UserIds []string `json:"userIds,omitempty"`
}

type UsersSetStatus struct {
	Status string `json:"status,omitempty"`
}

type UsersFollowReq struct {
	TargetId string          `json:"targetId,omitempty"`
	Action   SubscribeAction `json:"action,omitempty"`
}

type Auth interface {
	Login(context.Context, *EmailReq) (*AuthResp, error)
	Register(context.Context, *EmailReq) (*AuthResp, error)
	Vk(context.Context, *VkReq) (*AuthResp, error)
	CheckAuth(context.Context, *CheckAuthReq) (*AuthResp, error)
}

type CheckAuthReq struct {
	Token string `json:"token,omitempty"`
}

type EmailReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthResp struct {
	Token    string `json:"token,omitempty"`
	UserId   string `json:"userId,omitempty"`
	UserName string `json:"userName,omitempty"`
}

type VkReq struct {
	Code        string `json:"code,omitempty"`
	RedirectUrl string `json:"redirectUrl,omitempty"`
}

type Polls interface {
	Add(context.Context, *PollsAddReq) (*Poll, error)
	List(context.Context, *PollsListReq) (*PollsList, error)
	Vote(context.Context, *PollsVoteReq) (*Void, error)
	DeleteVote(context.Context, *PollsDeleteVoteReq) (*Void, error)
}

type PollsList struct {
	Items []*Poll `json:"items,omitempty"`
}

type PollsAddReq struct {
	Question string   `json:"question,omitempty"`
	Answers  []string `json:"answers,omitempty"`
}

type Poll struct {
	Id       string        `json:"id,omitempty"`
	Question string        `json:"question,omitempty"`
	Answers  []*PollAnswer `json:"answers,omitempty"`
}

type PollAnswer struct {
	Id         string `json:"id,omitempty"`
	Answer     string `json:"answer,omitempty"`
	VotedCount int    `json:"votedCount,omitempty"`
	IsVoted    bool   `json:"isVoted,omitempty"`
}

type PollsListReq struct {
	Ids []string `json:"ids,omitempty"`
}

type PollsVoteReq struct {
	PollId    string   `json:"pollId,omitempty"`
	AnswerIds []string `json:"answerIds,omitempty"`
}

type PollsDeleteVoteReq struct {
	PollId string `json:"pollId,omitempty"`
}
