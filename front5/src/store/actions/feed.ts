import {LoadingState, store} from "../store";
import * as types from "../../api/types";
import {parsePostsList} from "../helpers/posts";
import {api} from "../../api/api";

export function fetchFeed(type: "firstPage" | "nextPage") {
    if (store.getState().feed.stateFeed == LoadingState.FETCHING) {
        return;
    }

    if (type == "firstPage") {
        if (store.getState().feed.stateFeed != LoadingState.NONE) {
            return;
        }
    }

    store.dispatch({type: 'feed/setState', state: LoadingState.FETCHING});

    let cursor = '';
    if (store.getState().feed.pages.length > 0) {
        cursor = store.getState().feed.pages[store.getState().feed.pages.length - 1].nextCursor;
    }

    api("feed.list", {
        feedType: types.FeedType.FEED,
        cursor: cursor,
    } as types.FeedList).then((resp: types.PostsList) => {
        parsePostsList(resp);

        const postIds = resp.items.map(post => post.id);
        store.dispatch({type: 'feed/append', items: postIds, nextCursor: resp.nextCursor || ""});
        store.dispatch({type: 'feed/setState', state: LoadingState.DONE});
    })
}

export function fetchDiscover(type: "firstPage" | "nextPage") {
    if (store.getState().feed.stateDiscover == LoadingState.FETCHING) {
        return;
    }

    if (type == "firstPage") {
        if (store.getState().feed.stateDiscover != LoadingState.NONE) {
            return;
        }
    }

    store.dispatch({type: 'feed/setFeedDiscoverState', state: LoadingState.FETCHING});

    let cursor = '';
    if (store.getState().feed.discoverPages.length > 0) {
        cursor = store.getState().feed.discoverPages[store.getState().feed.discoverPages.length - 1].nextCursor;
    }

    api("feed.list", {
        feedType: types.FeedType.DISCOVER,
        cursor: cursor,
    } as types.FeedList).then((resp: types.PostsList) => {
        parsePostsList(resp);

        const postIds = resp.items.map(post => post.id);
        store.dispatch({type: 'feed/append', items: postIds, nextCursor: resp.nextCursor || ""});
        store.dispatch({type: 'feed/setFeedDiscoverState', state: LoadingState.DONE});
    })
}
