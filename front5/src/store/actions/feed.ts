import {LoadingState, store} from "../store";
import * as types from "../../api/types";
import {parsePostsList} from "../helpers/posts";
import {api} from "../../api/api";

export function fetchFeed() {
    const lockKey = "fetchFeed";
    if (store.getState().routing.fetchLockers[lockKey] === LoadingState.LOADING) {
        return;
    }
    store.dispatch({type: "routes/setFetchLocker", key: lockKey, state: LoadingState.LOADING});

    api("feed.list", {
        feedType: types.FeedType.DISCOVER,
        cursor: store.getState().feed.nextCursor
    } as types.FeedList).then((resp: types.PostsList) => {
        parsePostsList(resp);

        store.dispatch({type: "routes/setFetchLocker", key: lockKey, state: LoadingState.DONE});

        const postIds = resp.items.map(post => post.id);
        store.dispatch({type: "feed/append", items: postIds, nextCursor: resp.nextCursor || ""});
    })
}
