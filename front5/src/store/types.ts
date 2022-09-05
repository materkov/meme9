export type Composer = {}

export type Feed = [
    posts: Post[],
]

export type PostPage = [
    pagePost: string,
    post: Post,
]

export type UserPage = [
    user: User,
    posts: Post[],
]

export type User = {
    id: string;
    name?: string;
    href?: string;
}

export type Post = {
    id: string;
    fromId?: string;
    from?: User;
    text?: string;
    detailsURL?: string;
}

export type BrowseResult = {
    vkCallback?: VKCallback;

    componentName?: string;
    componentData?: any;
}

export type AddPostResponse = {
    post: Post;
}

export type VKCallback = {
    userId: string;
    authToken: string;
}
