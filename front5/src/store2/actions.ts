import {loadAPI} from "../store/fetcher";
import {Global, store} from "./store";
import * as types from "../store/types";
import {api} from "../store/types";
import {
    AppendFeed,
    AppendLikers,
    AppendPosts, DeleteFromFeed,
    PostLike,
    PostUnlike,
    SetLikes,
    SetOnline,
    SetPhoto,
    SetPost,
    SetRoute,
    SetUser,
    SetViewer
} from "./reducers";

export class Actions {
    prependFeed(post: types.Post): void {
        store.dispatch({type: "posts/set", post: post} as SetPost)
        store.dispatch({type: "feed/append", items: [post.id], prepend: true} as AppendFeed)
    }

    deletePost(postId: string) {
        store.dispatch({type: "feed/delete", postId: postId} as DeleteFromFeed)
    }

    loadFeed(): Promise<undefined> {
        const st = store.getState() as Global;
        if (st.feed.isLoaded) {
            return new Promise((resolve) => resolve);
        }

        return new Promise((resolve, reject) => {
            loadAPI(["/feed?feedType=DISCOVER&cursor="]).then(result => {
                for (let item of result) {
                    if (item.url.match("^/feed")) {
                        store.dispatch({type: "feed/append", items: item.items} as AppendFeed)
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

    loadViewer(): Promise<undefined> {
        return new Promise((resolve, reject) => {
            loadAPI(["/viewer"]).then(result => {
                for (let item of result) {
                    if (item.url == "/viewer") {
                        store.dispatch({type: "viewer/set", userId: item.viewerId || ""} as SetViewer)
                    }
                    if (item.url.match("^/users/\\d+$")) {
                        store.dispatch({type: "users/set", user: item} as SetUser)
                    }
                }
                resolve(undefined);
            })
        })
    }

    loadLikers(postId: string): Promise<undefined> {
        return new Promise((resolve, reject) => {
            loadAPI(["/posts/" + postId + "/liked?count=10"]).then(result => {
                for (let item of result) {
                    if (item.url.startsWith("/users/")) {
                        store.dispatch({type: "users/set", user: item} as SetUser)
                    }
                    if (item.url.startsWith("/posts/")) {
                        const parts = item.url.split("/");
                        store.dispatch({
                            type: "posts/appendLikers",
                            users: item.items || [],
                            postId: parts[2]
                        } as AppendLikers)
                    }
                }

                resolve(undefined);
            });
        })
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

    public postLike(id: string) {
        store.dispatch({type: 'posts/like', postId: id} as PostLike);

        api("/postLike", {id: id}).then(data => {
        })
    }

    public postUnlike(id: string) {
        store.dispatch({type: 'posts/unlike', postId: id} as PostUnlike);

        api("/postUnlike", {id: id}).then(data => {
        })
    }

    public setRoute(url: string) {
        store.dispatch({type: 'routes/set', url: url} as SetRoute);
    }
}


export const actions = new Actions();
