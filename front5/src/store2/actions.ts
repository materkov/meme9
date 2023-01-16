import {loadAPI} from "../store/fetcher";
import {Global, store} from "./store";
import {AppendFeed, AppendPosts, SetIsFollowing, SetLikes, SetOnline, SetPhoto, SetPost, SetUser} from "./reducers";

export class Actions {
    loadFeed(): Promise<undefined> {
        const st = store.getState() as Global;
        if (st.feed.isLoaded) {
            return new Promise((resolve) => resolve);
        }

        return new Promise((resolve, reject) => {
            loadAPI(["/feed?feedType=DISCOVER&cursor="]).then(result => {
                for (let item of result) {
                    if (item.url.match("^/feed")) {
                        store.dispatch({type: "feed/append", items: item.items || []} as AppendFeed)
                    }
                    if (item.url.match("^/posts/\\d+$")) {
                        store.dispatch({type: "posts/set", post: item} as SetPost)
                    }
                    if (item.url.match("^/posts/\\d+/liked")) {
                        const parts = item.url.split("/");

                        store.dispatch({
                            type: "posts/setLikes",
                            postId: parts[2],
                            isLiked: item.isViewerLiked || false,
                            count: item.totalCount || 0
                        } as SetLikes)
                    }
                    if (item.url.match("^/photos/\\d+$")) {
                        store.dispatch({type: "photos/set", photo: item} as SetPhoto)
                    }
                    if (item.url.match("^/users/\\d+/online$")) {
                        const parts = item.url.split("/")
                        store.dispatch({type: "online/set", online: item, userId: parts[2]} as SetOnline)
                    }
                    if (item.url.match("^/users/\\d+$")) {
                        store.dispatch({type: "users/set", user: item} as SetUser)
                    }
                }

                resolve(undefined);
            });
        });
    }

    loadUserPage(userId: string): Promise<undefined> {
        return new Promise((resolve, reject) => {
            loadAPI([
                "/users/" + userId,
                `/users/${userId}/posts?count=10`,
                "/users/" + userId + "/followers",
                "/users/" + userId + "/following",
            ]).then(result => {
                for (let item of result) {
                    if (item.url.match("^/posts/\\d+$")) {
                        store.dispatch({type: "posts/set", post: item} as SetPost)
                    }
                    if (item.url.match("^/posts/\\d+/liked")) {
                        const parts = item.url.split("/");

                        store.dispatch({
                            type: "posts/setLikes",
                            postId: parts[2],
                            isLiked: item.isViewerLiked || false,
                            count: item.totalCount || 0
                        } as SetLikes)
                    }
                    if (item.url.match("^/photos/\\d+$")) {
                        store.dispatch({type: "photos/set", photo: item} as SetPhoto)
                    }
                    if (item.url.match("^/users/\\d+/online$")) {
                        const parts = item.url.split("/")
                        store.dispatch({type: "online/set", online: item, userId: parts[2]} as SetOnline)
                    }
                    if (item.url.match("^/users/\\d+/followers")) {
                        const parts = item.url.split("/");
                        store.dispatch({
                            type: "users/setIsFollowing",
                            isFollowing: item.isFollowing || false,
                            userId: parts[2]
                        } as SetIsFollowing)
                    }
                    if (item.url.match("^/users/\\d+/posts")) {
                        const parts = item.url.split("/")
                        store.dispatch({
                            type: "users/appendPosts",
                            userId: parts[2],
                            posts: item.items || []
                        } as AppendPosts)

                    }
                    if (item.url.match("^/users/\\d+$")) {
                        store.dispatch({type: "users/set", user: item} as SetUser)
                    }
                }
                resolve(undefined);
            })
        })
    }
}

export const actions = new Actions();
