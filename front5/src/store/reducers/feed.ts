import {Global, LoadingState} from "../store";

export interface SetState {
    type: 'feed/setState'
    state: LoadingState
}

export function setLoadingState(state: Global, data: SetState): Global {
    return {
        ...state,
        feed: {
            ...state.feed,
            state: data.state,
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
