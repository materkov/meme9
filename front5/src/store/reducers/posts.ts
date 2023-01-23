import {Global} from "../store";
import * as types from "../../api/types";

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
