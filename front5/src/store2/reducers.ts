import {Global} from "./store";
import {PostLike} from "../components/PostLike";
import * as types from "../store/types";

export interface PostLike {
    type: 'posts/like';
    postId: string;
}

export function postLikeReducer(state: Global, data: PostLike) {
    return {
        ...state,
        posts: {
            ...state.posts,
            isLiked: {
                ...state.posts.isLiked,
                [data.postId]: true,
            },
            likesCount: {
                ...state.posts.likesCount,
                [data.postId]: (state.posts.likesCount[data.postId] || 0) + 1,
            }
        }
    }
}

export interface PostUnlike {
    type: 'posts/unlike';
    postId: string;
}

export function postUnlike(state: Global, data: PostUnlike) {
    const isLiked = {...state.posts.isLiked};
    delete isLiked[data.postId];

    const likesCount = {...state.posts.likesCount};
    const count = likesCount[data.postId];
    if (count > 1) {
        likesCount[data.postId]--;
    } else {
        delete likesCount[data.postId];
    }

    return {
        ...state,
        posts: {
            ...state.posts,
            isLiked: isLiked,
            likesCount: likesCount
        }
    }
}

export interface SetRoute {
    type: 'routes/set'
    url: string
}

export function setRouteReducer(state: Global, data: SetRoute): Global {
    return {
        ...state,
        routing: {
            ...state.routing,
            url: data.url,
        }
    }
}

export interface SetPost {
    type: 'posts/set'
    post: types.Post
}

export function setPost(state: Global, data: SetPost): Global {
    return {
        ...state,
        posts: {
            ...state.posts,
            byId: {
                ...state.posts.byId,
                [data.post.id]: data.post,
            }
        }
    }
}

export interface SetUser {
    type: 'users/set'
    user: types.User
}

export function setUser(state: Global, data: SetUser): Global {
    return {
        ...state,
        users: {
            ...state.users,
            byId: {
                ...state.users.byId,
                [data.user.id]: data.user,
            }
        }
    }
}

export interface SetOnline {
    type: 'online/set'
    userId: string;
    online: types.Online
}

export function setOnline(state: Global, data: SetOnline): Global {
    return {
        ...state,
        online: {
            ...state.online,
            byId: {
                ...state.online.byId,
                [data.userId]: data.online,
            }
        }
    }
}

export interface SetPhoto {
    type: 'photos/set'
    photo: types.Photo;
}

export function setPhoto(state: Global, data: SetPhoto): Global {
    return {
        ...state,
        photos: {
            ...state.photos,
            byId: {
                ...state.photos.byId,
                [data.photo.id]: data.photo,
            }
        }
    }
}

export interface AppendFeed {
    type: 'feed/append'
    items: string[]
    prepend?: boolean
}

export function appendFeed(state: Global, data: AppendFeed): Global {
    return {
        ...state,
        feed: {
            ...state.feed,
            isLoaded: true,
            items: data.prepend ? [...data.items, ...state.feed.items] : [...state.feed.items, ...data.items],
        }
    }
}

export interface DeleteFromFeed {
    type: 'feed/delete'
    postId: string
}

export function deleteFromFeed(state: Global, data: DeleteFromFeed): Global {
    return {
        ...state,
        feed: {
            ...state.feed,
            items: state.feed.items.filter(item => item != data.postId)
        }
    }
}

export interface SetLikes {
    type: 'posts/setLikes'
    postId: string
    count: number
    isLiked: boolean
}

export function setLikes(state: Global, data: SetLikes): Global {
    const likesCount = {...state.posts.likesCount};
    if (data.count > 0) {
        likesCount[data.postId] = data.count
    } else {
        delete likesCount[data.postId];
    }

    const isLiked = {...state.posts.isLiked};
    if (data.isLiked) {
        isLiked[data.postId] = true;
    } else {
        delete isLiked[data.postId];
    }

    return {
        ...state,
        posts: {
            ...state.posts,
            likesCount: likesCount,
            isLiked: isLiked
        }
    }
}


export interface AppendLikers {
    type: 'posts/appendLikers'
    postId: string
    users: string[]
}

export function appendLikers(state: Global, data: AppendLikers): Global {
    return {
        ...state,
        posts: {
            ...state.posts,
            likers: {
                ...state.posts.likers,
                [data.postId]: [...(state.posts.likers[data.postId] || []), ...data.users],
            }
        }
    }
}

export interface AppendPosts {
    type: 'users/appendPosts'
    userId: string
    posts: string[]
}

export function appendPosts(state: Global, data: AppendPosts): Global {
    return {
        ...state,
        users: {
            ...state.users,
            posts: {
                ...state.users.posts,
                [data.userId]: [...(state.users.posts[data.userId] || []), ...data.posts],
            }
        }
    }
}


export interface SetViewer {
    type: 'viewer/set'
    userId: string
}

export function setViewer(state: Global, data: SetViewer): Global {
    return {
        ...state,
        viewer: {
            ...state.viewer,
            isLoaded: true,
            id: data.userId,
        }
    }
}

export interface SetIsFollowing {
    type: 'users/setIsFollowing'
    userId: string;
    isFollowing: boolean;
}

export function setIsFollowing(state: Global, data: SetIsFollowing): Global {
    const isFollowing = {...state.users.isFollowing};
    if (data.isFollowing) {
        isFollowing[data.userId] = true;
    } else {
        delete isFollowing[data.userId];
    }

    return {
        ...state,
        users: {
            ...state.users,
            isFollowing: isFollowing,
        }
    }
}

export interface SetToken {
    type: 'auth/setToken';
    token: string;
}

export function setToken(state: Global, data: SetToken): Global {
    return {
        ...state,
        routing: {
            ...state.routing,
            accessToken: data.token,
        }
    }
}

export type AnyAction = PostLike | PostUnlike | SetRoute | SetPost | SetUser | SetOnline | AppendFeed
    | SetPhoto | SetLikes | AppendLikers | AppendPosts | SetViewer | DeleteFromFeed | SetIsFollowing
    | SetToken
    ;
