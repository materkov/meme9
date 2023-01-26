export type User = {
    id: string;
    name?: string;
    avatar?: string;
    bio?: string;
    online?: Online;
}

export type Post = {
    id: string;
    userId: string;
    user?: User;
    date: string;
    text: string;
    canDelete?: boolean;
    isDeleted?: boolean;
    photoId?: string;
    photo?: Photo;

    likesConnection?: PostsLikesConnection;
}

export type Photo = {
    id: string;
    url: string;
    address: string;
    width: number;
    height: number;
    thumbs: PhotoThumb[];
}

export type PhotoThumb = {
    width: number;
    height: number;
    address: string;
}


export interface PostLikeData {
    isViewerLiked: boolean
    totalCount: number
    items: string
}

export interface Edges {
    totalCount: number
    items: string[]
    nextCursor: string
}

export interface FollowersEdges {
    totalCount: number
    items: string[]
    nextCursor: string
    isFollowing: boolean
}

export interface Viewer {
    url: string;
    viewerId: string;
}

export interface Online {
    url: string;
    isOnline: boolean;
}

export interface PostsAdd {
    method: 'posts.add';
    text: string;
    photo: string;
}

export interface PostsDelete {
    method: 'posts.delete';
    id: string;
}

export interface PostsLike {
    method: 'posts.like';
    postId: string;
}

export interface PostsUnlike {
    method: 'posts.unlike';
    postId: string;
}

export interface UsersFollow {
    method: 'users.follow';
    userId: string;
}

export interface UsersUnfollow {
    method: 'users.unfollow';
    userId: string;
}

export interface UsersEdit {
    method: 'users.edit';
    userId: string;
    name: string;
}

export interface AuthVkCallback {
    method: 'auth.vkCallback';
    code: string;
    redirectUri: string;
}

export interface Authorization {
    token: string;
    user: User;
}

export interface AuthEmailLogin {
    method: 'auth.emailLogin';
    email: string;
    password: string;
}

export interface AuthEmailRegister {
    method: 'auth.emailRegister';
    email: string;
    password: string;
}

export interface UsersSetAvatar {
    method: 'users.setAvatar';
    uploadToken: string;
}

export interface PostsGetLikesConnection {
    method: 'posts.getLikes';
    postId: string
    count: number
}

export interface PostsLikesConnection {
    totalCount: number
    isViewerLiked: boolean

    items: User[]
}

export enum FeedType {
    FEED = 'FEED',
    DISCOVER = 'DISCOVER',
}

export interface FeedList {
    method: 'feed.list';
    feedType: FeedType
}

export interface PostsList {
    items: Post[]
    totalCount: number;
    nextCursor: string
}

export interface UsersPostsList {
    method: 'users.posts.list';
    userId: string;
    count: number;
}

export interface UserFollowingCount {
    userId: string;
}

export interface UserFollowersCount {
    userId: string;
}

export type ApiRequest = PostsAdd | PostsDelete | PostsLike | PostsUnlike |
    UsersFollow | UsersUnfollow | UsersEdit | UsersSetAvatar |
    AuthVkCallback | AuthEmailLogin | AuthEmailRegister |
    PostsGetLikesConnection | FeedList | UsersPostsList
    ;