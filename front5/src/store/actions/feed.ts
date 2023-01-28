import {store} from "../store";
import * as types from "../../api/types";
import {parsePostsList} from "../helpers/posts";
import {api} from "../../api/api";

export function fetchFeed() {
    let cursor = '';
    if (store.getState().feed.pages.length > 0) {
        cursor = store.getState().feed.pages[store.getState().feed.pages.length - 1].nextCursor;
    }

    api("feed.list", {
        feedType: types.FeedType.DISCOVER,
        cursor: cursor,
    } as types.FeedList).then((resp: types.PostsList) => {
        parsePostsList(resp);

        const postIds = resp.items.map(post => post.id);
        store.dispatch({type: 'feed/append', items: postIds, nextCursor: resp.nextCursor || ""});
    })
}
