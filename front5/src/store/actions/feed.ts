import {store} from "../store";
import * as types from "../../api/types";
import {parsePostsList} from "../helpers/posts";
import {AppendFeed} from "../reducers/feed";
import {api} from "../../api/api";

export function loadFeed() {
    const st = store.getState();
    if (st.feed.isLoaded) {
        return new Promise((resolve) => resolve);
    }

    api("feed.list", {feedType: types.FeedType.DISCOVER} as types.FeedList).then((resp: types.PostsList) => {
        parsePostsList(resp);

        const postIds = resp.items.map(post => post.id);
        store.dispatch({type: "feed/append", items: postIds} as AppendFeed);

    })
}
