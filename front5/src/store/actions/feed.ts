import {LoadingState, store} from "../store";
import * as types from "../../api/types";
import {parsePostsList} from "../helpers/posts";
import {api} from "../../api/api";

export function fetchFeed() {
    const st = store.getState();
    if (st.feed.state !== LoadingState.NONE) {
        return;
    }

    store.dispatch({type: "feed/setState", state: LoadingState.LOADING});

    api("feed.list", {feedType: types.FeedType.DISCOVER} as types.FeedList).then((resp: types.PostsList) => {
        parsePostsList(resp);

        store.dispatch({type: "feed/setState", state: LoadingState.DONE});

        const postIds = resp.items.map(post => post.id);
        store.dispatch({type: "feed/append", items: postIds});
    })
}
