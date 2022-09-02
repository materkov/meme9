export type Post = {
    id: string;
    text?: string;
    date?: number;
    user?: User;
}

export type PostsListRequest = {
    id: string[];
    fields: string;
}

export type PostsListResponse = {
    items: Post[];
}

export const PostsListRequest = (req: PostsListRequest): Operation => ({
    method: "posts.list",
    params: req,
})

export type FeedListRequest = {
    fields: string;
}

export type FeedListResponse = {
    items: Post[];
}

export const FeedList = (req: FeedListRequest): Operation => ({
    method: "feed.list",
    params: req,
})

export type User = {
    id: string;
    name?: string;
}

export type UsersListRequest = {
    id: string[];
    fields: string;
}

export type UsersListResponse = {
    items: User[];
}

export const UsersList = (req: UsersListRequest): Operation => ({
    method: "users.list",
    params: req,
})

type Operation = {
    method: string;
    params: any;
}
