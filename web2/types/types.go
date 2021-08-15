package types

type UniversalRenderer struct {
	NewPostRenderer  *NewPostRenderer  `json:"newPostRenderer,omitempty"`
	PostRenderer     *PostRenderer     `json:"postRenderer,omitempty"`
	UserPageRenderer *UserPageRenderer `json:"userPageRenderer,omitempty"`
	VkAuthRenderer   *VkAuthRenderer   `json:"vkAuthRenderer,omitempty"`
	FeedRenderer     *FeedRenderer     `json:"feedRenderer,omitempty"`
}

type NewPostRenderer struct {
	SendLabel string `json:"sendLabel,omitempty"`
}

type FeedRenderer struct {
	Posts []*PostRenderer `json:"posts,omitempty"`
}

type PostRenderer struct {
	ID         string `json:"id,omitempty"`
	AuthorName string `json:"authorName,omitempty"`
	AuthorHref string `json:"authorHref,omitempty"`
	Text       string `json:"text,omitempty"`
}

type UserPageRenderer struct {
	UserName string          `json:"userName,omitempty"`
	UserID   string          `json:"userId,omitempty"`
	Posts    []*PostRenderer `json:"posts,omitempty"`
}

type VkAuthRenderer struct {
	URL string `json:"url,omitempty"`
}
