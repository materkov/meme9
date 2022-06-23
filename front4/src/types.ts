export type Query = {
    feed?: Post[];
    vkAuthUrl?: string;
    mutation?: Mutation;
    viewer?: User;
    node?: Node;
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
}

export type UserParams = {
    name?: {};
    avatar?: {};
    posts?: UserPostsConnectionFields;
}

export type UserPostsConnectionFieldsEdges = {
    inner?: PostParams;
}

export type UserPostsConnection = {
    id: string;
    type: "UserPostsConnection";

    totalCount?: number;
    edges?: Post[];
}

export type UserPostsConnectionFields = {
    totalCount?: {};
    edges?: UserPostsConnectionFieldsEdges;
};

export type Post = {
    type: "Post";
    id: string;
    text?: string;
    user?: User;
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
}

export type MutationAddPost = {
    text?: string;
}

export type MutationVKAuthCallback = {
    url?: string;
}