import {Global, store} from "../store";
import {AppendFeed, SetLikes, SetOnline, SetPhoto, SetPost, SetUser} from "../reducers";
import {api2, FeedList, FeedType, PostsList} from "../../store/types";
import {parsePostsList} from "../helpers/posts";

export function loadFeed(): Promise<undefined> {
    const st = store.getState() as Global;
    if (st.feed.isLoaded) {
        return new Promise((resolve) => resolve);
    }

    return new Promise((resolve, reject) => {
        api2("feed.list", {feedType: FeedType.DISCOVER} as FeedList).then((resp: PostsList) => {
            parsePostsList(resp);

            const postIds = resp.items.map(post => post.id);
            store.dispatch({type: "feed/append", items: postIds} as AppendFeed);

            resolve(undefined);
        })
    });
}
