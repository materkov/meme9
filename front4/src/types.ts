export type Query = {
    type: "Query";
    id: string;
    feed?: string;
    vkAuthUrl?: string;
    mutation?: Mutation;
    viewer?: string;
    node?: Node;
}

export type Feed = {
    type: "Feed";
    id: string;
    feed: string[];
}

export type Mutation = {
    vkAuth?: VKAuth;
}

export type VKAuth = {
    token?: string;
}

export type QueryParams = {
    feed?: QueryFeed;
    mutation?: QueryMutation;
    vkAuthUrl?: {};
    viewer?: QueryViewer;
    node?: QueryNode;
}

export type QueryNode = {
    id?: string;
    inner?: NodeParams;
}

export type NodeParams = {
    onPost?: PostParams;
    onUser?: UserParams;
}

export type Node = Post | User;

export type QueryViewer = {
    inner?: UserParams;
}

export type QueryFeed = {
    userId?: number;
    inner?: PostParams;
}

export type QueryMutation = {
    inner?: MutationParams;
}

export type User = {
    type: "User";
    id: string;
    name?: string;
    avatar?: string;
    posts?: UserPostsConnection;
    isFollowing?: boolean;
}

export type UserParams = {
    name?: {};
    avatar?: {};
    posts?: UserPostsConnectionFields;
    isFollowing?: {};
}

export type UserPostsConnectionFieldsEdges = {
    inner?: PostParams;
}

export type UserPostsConnection = {
    id: string;
    type: "UserPostsConnection";

    totalCount?: number;
    edges?: string[];
}

export type UserPostsConnectionFields = {
    totalCount?: {};
    edges?: UserPostsConnectionFieldsEdges;
};

export type Post = {
    type: "Post";
    id: string;
    text?: string;
    user?: string;
    date?: number;
}

export type PostParams = {
    date?: {};
    text?: PostText;
    user?: PostUser;
}

export type PostText = {
    maxLength?: number;
}

export type PostUser = {
    inner?: UserParams;
}

export type MutationParams = {
    addPost?: MutationAddPost;
    vkAuthCallback?: MutationVKAuthCallback;
    follow?: MutationFollow;
    unfollow?: MutationUnfollow;
}

export type MutationFollow = {
    userId?: string;
}

export type MutationUnfollow = {
    userId?: string;
}

export type MutationAddPost = {
    text?: string;
}

export type MutationVKAuthCallback = {
    url?: string;
}

export type CurrentRoute = {
    id: string;
    type: "CurrentRoute";
    url: string;
}

export type Viewer = {
    id: string;
    type: "Viewer";
    userId?: string;
}

export type AuthToken = {
    id: string;
    type: "AuthToken";
    token?: string;
}

export type VkAuthURL = {
    id: string;
    type: "VkAuthURL";
    url?: string;
}
