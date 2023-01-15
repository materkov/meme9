import {createStore} from "redux";
import * as types from "../store/types";
import {
    AnyAction,
    appendFeed,
    appendLikers,
    appendPosts, deleteFromFeed,
    postLikeReducer,
    postUnlike,
    setLikes,
    setOnline,
    setPhoto,
    setPost,
    setRouteReducer,
    setUser,
    setViewer
} from "./reducers";

export interface Global {
    routing: {
        url: string
    }

    posts: {
        byId: { [id: string]: types.Post }
        likesCount: { [id: string]: number }
        likers: { [id: string]: string[] }
        isLiked: { [id: string]: boolean }
    }

    feed: {
        isLoading: boolean
        isLoaded: boolean
        items: string[]
    }

    users: {
        byId: { [id: string]: types.User }

        postsCount: { [id: string]: number }
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
        url: location.pathname + location.search,
    },
    posts: {
        byId: {},
        likesCount: {},
        likers: {},
        isLiked: {},
    },
    feed: {
        isLoaded: false,
        items: [],
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

export const store = createStore<Global, AnyAction, any, any>((state = global, action: AnyAction) => {
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
        default:
            return state
    }
}, undefined, window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__())

// @ts-ignore
window.global = store;
