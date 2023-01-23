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
    text: string;
    photo: string;
}

export interface PostsDelete {
    id: string;
}

export interface PostsLike {
    postId: string;
}

export interface PostsUnlike {
    postId: string;
}

export interface UsersFollow {
    userId: string;
}

export interface UsersUnfollow {
    userId: string;
}

export interface UsersEdit {
    userId: string;
    name: string;
}

export interface AuthVkCallback {
    code: string;
    redirectUri: string;
}

export interface Authorization {
    token: string;
    user: User;
}

export interface AuthEmailLogin {
    email: string;
    password: string;
}

export interface AuthEmailRegister {
    email: string;
    password: string;
}

export interface UsersSetAvatar {
    uploadToken: string;
}

export interface PostsGetLikesConnection {
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
    feedType: FeedType
}

export interface PostsList {
    items: Post[]
    totalCount: number;
    nextCursor: string
}

export interface UsersPostsList {
    userId: string;
    count: number;
}

export interface UserFollowingCount {
    userId: string;
}

export interface UserFollowersCount {
    userId: string;
}