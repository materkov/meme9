export interface Query {
    feed: Post[];
}

export interface QueryParams {
    feed?: QueryFeed;
    mutation?: QueryMutation;
}

export interface QueryFeed {
    include?: boolean;
    userId?: number;
    inner?: PostParams;
}

export interface QueryMutation {
    include?: boolean;
    inner?: MutationParams;
}

export interface User {
    type: "User";
    id: string;
    name?: string;
}

export interface UserParams {
    name?: SimpleParams;
}

export interface SimpleParams {
    include?: boolean;
}

export interface Post {
    type: "Post";
    id: string;
    text?: string;
    user?: User;
    date?: number;
}

export interface PostParams {
    date?: SimpleParams;
    text?: PostText;
    user?: PostUser;
}

export interface PostText {
    include?: boolean;
    maxLength?: number;
}

export interface PostUser {
    include?: boolean;
    inner?: UserParams;
}

export interface MutationParams {
    addPost?: MutationAddPost;
}

export interface MutationAddPost {
    include?: boolean;
    text?: string;
}
