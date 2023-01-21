import {Global, store} from "../store2/store";

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

export function api2(method: string, args: any): Promise<any> {
    const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me/api2?method=" : "http://localhost:8000/api2?method=";
    const state = store.getState() as Global;

    return new Promise((resolve, reject) => {
        return fetch(apiHost + method, {
            method: 'POST',
            headers: {
                'authorization': 'Bearer ' + state.routing.accessToken,
            },
            body: JSON.stringify(args),
        })
            .then(resp => resp.json())
            .then(resp => {
                if (resp.error) {
                    reject(resp.error);
                } else {
                    resolve(resp.data);
                }
            })
    })
}

export function uploadApi(file: File): Promise<string> {
    const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me" : "http://localhost:8000";

    const form = new FormData();
    form.set("file", file);

    return new Promise((resolve, reject) => {
        return fetch(apiHost + "/upload", {
            method: 'POST',
            headers: {},
            body: form,
        })
            .then(resp => {
                if (resp.status == 200) {
                    resp.json().then((res) => resolve(res.uploadToken));
                } else {
                    resp.text().then(reject);
                }
            })
    })
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
