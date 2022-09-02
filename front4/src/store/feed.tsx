import {QueryParams} from "../types";
import {PostQuery} from "../components/Post";
import {api} from "../api";
import {storeOnChanged} from "./store";

export function storeFeedLoad() {
    const feedQuery: QueryParams = {
        feed: {
            inner: PostQuery,
        },
        viewer: {
            inner: {},
        }
    }

    api(feedQuery).then(data => {
        // TODO lineraize
        // @ts-ignore
        window.store["feed:fake:id"] = {
            type: "Feed",
            feed: [],
        };

        // @ts-ignore
        window.store["query"] = {
            type: "Query",
            feed: "feed:fake:id",
            // @ts-ignore
            viewer: data.viewer?.id || 0,
        };
        // @ts-ignore
        if (data.viewer?.id) {
            // @ts-ignore
            window.store[data.viewer.id] = data.viewer.id;
        }

        // @ts-ignore
        for (let post of data.feed || []) {
            // @ts-ignore
            window.store[post.user.id] = post.user;
            // @ts-ignore
            post.user = post.user.id;
            // @ts-ignore
            window.store[post.id] = post;

            // @ts-ignore
            window.store["feed:fake:id"].feed.push(post.id);
        }

        storeOnChanged();
    })
}