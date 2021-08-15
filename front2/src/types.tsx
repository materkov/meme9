export interface UniversalRenderer {
    newPostRenderer?: NewPostRenderer
    postRenderer?: PostRenderer
    userPageRenderer?: UserPageRenderer
    vkAuthRenderer?: VkAuthRenderer
    feedRenderer?: FeedRenderer
}

export interface NewPostRenderer {
    sendLabel?: string
}

export interface FeedRenderer {
    posts?: PostRenderer[]
}

export interface PostRenderer {
    id?: string;
    authorName?: string
    authorHref?: string
    text?: string
}

export interface UserPageRenderer {
    userName?: string
    userId?: string
    posts?: PostRenderer[]
}

export interface VkAuthRenderer {
    url?: string
}
