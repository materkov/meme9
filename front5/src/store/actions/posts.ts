import * as types from "../../api/types";
import {store} from "../store";
import {api} from "../../api/api";

export function deletePost(postId: string) {
    api("posts.delete", {
        id: postId,
    } as types.PostsDelete);

    store.dispatch({type: "feed/delete", postId: postId})
}

export function addPost(data: types.PostsAdd): Promise<void> {
    return new Promise((resolve, reject) => {
        api("posts.add", data).then((post: types.Post) => {
            store.dispatch({type: "posts/set", post: post})
            store.dispatch({type: "feed/append", items: [post.id], prepend: true, nextCursor: ''})

            resolve();
        }).catch(() => {
            reject();
        })
    })
}

export function like(postId: string) {
    const state = store.getState();

    store.dispatch({
        type: 'posts/setLikes',
        postId: postId,
        count: (state.posts.likesCount[postId] || 0) + 1,
        isLiked: true,
    });

    api("posts.like", {postId: postId} as types.PostsLike).then((resp: types.PostLikeData) => {
        store.dispatch({
            type: 'posts/setLikes',
            postId: postId,
            count: resp.totalCount || 0,
            isLiked: true,
        });
    });
}

export function unlike(postId: string) {
    const state = store.getState();

    store.dispatch({
        type: 'posts/setLikes',
        postId: postId,
        count: (state.posts.likesCount[postId] || 0) - 1,
        isLiked: false,
    });

    api("posts.unlike", {postId: postId} as types.PostsUnlike).then((resp: types.PostLikeData) => {
        store.dispatch({
            type: 'posts/setLikes',
            postId: postId,
            count: resp.totalCount || 0,
            isLiked: false,
        });
    });
}

export function loadLikers(postId: string): Promise<void> {
    return new Promise((resolve) => {
        api("posts.getLikesConnection", {
            postId: postId,
            count: 10
        } as types.PostsGetLikesConnection).then((resp: types.PostsLikesConnection) => {
            for (let user of resp.items) {
                store.dispatch({type: "users/set", user: user})
            }

            const likers: string[] = resp.items.map(u => u.id);
            store.dispatch({
                type: "posts/appendLikers",
                users: likers,
                postId: postId
            });

            resolve();
        })
    })
}
