import {Global} from "../store";
import * as types from "../../api/types";

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

export interface SetPostsCount {
    type: 'users/setPostsCount'
    userId: string
    count: number | null
}

export function setPostsCount(state: Global, data: SetPostsCount): Global {
    return {
        ...state,
        users: {
            ...state.users,
            postsCount: {
                ...state.users.postsCount,
                [data.userId]: data.count,
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

export interface SetFollowingCount {
    type: 'users/setFollowingCount'
    userId: string;
    count: number;
}

export function setFollowingCount(state: Global, data: SetFollowingCount): Global {
    return {
        ...state,
        users: {
            ...state.users,
            followingCount: {
                ...state.users.followingCount,
                [data.userId]: data.count,
            }
        }
    }
}

export interface SetFollowersCount {
    type: 'users/setFollowersCount'
    userId: string;
    count: number;
    isViewerFollowing: boolean;
}

export function setFollowersCount(state: Global, data: SetFollowersCount): Global {
    return {
        ...state,
        users: {
            ...state.users,
            followersCount: {
                ...state.users.followersCount,
                [data.userId]: data.count,
            },
            isFollowing: {
                ...state.users.isFollowing,
                [data.userId]: data.isViewerFollowing
            }
        }
    }
}