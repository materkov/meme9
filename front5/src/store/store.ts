import {createStore} from "redux";
import * as types from "../api/types";
import {
    AppendLikers,
    appendLikers,
    PostLike,
    postLikeReducer,
    PostUnlike,
    postUnlike,
    SetLikes,
    setLikes,
    SetPost,
    setPost
} from "./reducers/posts";
import {
    AppendPosts,
    appendPosts,
    SetFollowersCount,
    setFollowersCount,
    SetFollowingCount,
    setFollowingCount,
    SetIsFollowing,
    setIsFollowing,
    SetPostsCount,
    setPostsCount,
    SetUser,
    setUser
} from "./reducers/users";
import {SetRoute, setRouteReducer} from "./reducers/routes";
import {SetViewer, setViewer} from "./reducers/viewer";
import {AppendFeed, appendFeed, DeleteFromFeed, deleteFromFeed} from "./reducers/feed";
import {SetOnline, setOnline} from "./reducers/online";
import {SetPhoto, setPhoto} from "./reducers/photos";
import {SetToken, setToken} from "./reducers/auth";

export interface Page {
    items: string[]
    nextCursor: string
}

export interface Global {
    routing: {
        accessToken: string
        url: string
    }

    posts: {
        byId: { [id: string]: types.Post }
        likesCount: { [id: string]: number }
        likers: { [id: string]: string[] }
        isLiked: { [id: string]: boolean }
    }

    feed: {
        pages: Page[]
    }

    users: {
        byId: { [id: string]: types.User }

        postsCount: { [id: string]: number | null }
        followersCount: { [id: string]: number }
        followingCount: { [id: string]: number }

        isFollowing: { [id: string]: boolean }
        posts: { [id: string]: string[] }
    }

    online: {
        byId: { [id: string]: types.Online }
    }

    photos: {
        byId: { [id: string]: types.Photo }
    }

    viewer: {
        isLoaded: boolean;
        id: string;
    }
}

const global: Global = {
    routing: {
        accessToken: localStorage.getItem("authToken") || "",
        url: location.pathname + location.search,
    },
    posts: {
        byId: {},
        likesCount: {},
        likers: {},
        isLiked: {},
    },
    feed: {
        pages: [],
    },
    users: {
        byId: {},

        posts: {},
        isFollowing: {},
        followingCount: {},
        followersCount: {},
        postsCount: {},
    },
    online: {
        byId: {},
    },
    photos: {
        byId: {},
    },
    viewer: {
        isLoaded: false,
        id: "",
    }
}

export type MyAnyAction = PostLike | PostUnlike | SetRoute | SetPost | SetUser | SetOnline | AppendFeed
    | SetPhoto | SetLikes | AppendLikers | AppendPosts | SetViewer | DeleteFromFeed | SetIsFollowing
    | SetToken | SetPostsCount | SetFollowingCount | SetFollowersCount
    ;


export const store = createStore<Global, MyAnyAction, void, void>((state = global, action: MyAnyAction) => {
    switch (action.type) {
        case 'posts/like':
            return postLikeReducer(state, action)
        case 'posts/unlike':
            return postUnlike(state, action)
        case 'routes/set':
            return setRouteReducer(state, action)
        case 'posts/set':
            return setPost(state, action)
        case 'users/set':
            return setUser(state, action)
        case 'online/set':
            return setOnline(state, action)
        case 'photos/set':
            return setPhoto(state, action)
        case 'feed/append':
            return appendFeed(state, action)
        case 'posts/setLikes':
            return setLikes(state, action)
        case 'posts/appendLikers':
            return appendLikers(state, action)
        case 'users/appendPosts':
            return appendPosts(state, action)
        case 'viewer/set':
            return setViewer(state, action)
        case 'feed/delete':
            return deleteFromFeed(state, action)
        case 'users/setIsFollowing':
            return setIsFollowing(state, action)
        case 'users/setPostsCount':
            return setPostsCount(state, action)
        case 'users/setFollowersCount':
            return setFollowersCount(state, action)
        case 'users/setFollowingCount':
            return setFollowingCount(state, action)
        case 'auth/setToken':
            return setToken(state, action)
        default:
            return state
    }
// @ts-ignore
}, undefined, window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__())

// @ts-ignore
window.global = store;
