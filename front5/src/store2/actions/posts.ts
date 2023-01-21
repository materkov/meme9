import * as types from "../../api/types";
import {store} from "../store";
import {AppendFeed, DeleteFromFeed} from "../reducers/feed";
import {AppendLikers, SetLikes, SetPost} from "../reducers/posts";
import {SetUser} from "../reducers/users";
import {api} from "../../api/api";

export function deletePost(postId: string) {
    api("posts.delete", {
        id: postId,
    } as types.PostsDelete);

    store.dispatch({type: "feed/delete", postId: postId} as DeleteFromFeed)
}

export function add(data: types.PostsAdd): Promise<void> {
    return new Promise((resolve, reject) => {
        api("posts.add", data).then((post: types.Post) => {
            store.dispatch({type: "posts/set", post: post} as SetPost)
            store.dispatch({type: "feed/append", items: [post.id], prepend: true} as AppendFeed)

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
    } as SetLikes);

    api("posts.like", {postId: postId} as types.PostsLike).then((resp: types.PostLikeData) => {
        store.dispatch({
            type: 'posts/setLikes',
            postId: postId,
            count: resp.totalCount || 0,
            isLiked: true,
        } as SetLikes);
    });
}

export function unlike(postId: string) {
    const state = store.getState();

    store.dispatch({
        type: 'posts/setLikes',
        postId: postId,
        count: (state.posts.likesCount[postId] || 0) - 1,
        isLiked: false,
    } as SetLikes);

    api("posts.unlike", {postId: postId} as types.PostsUnlike).then((resp: types.PostLikeData) => {
        store.dispatch({
            type: 'posts/setLikes',
            postId: postId,
            count: resp.totalCount || 0,
            isLiked: false,
        } as SetLikes);
    });
}

export function loadLikers(postId: string): Promise<void> {
    return new Promise((resolve, reject) => {
        api("posts.getLikesConnection", {
            postId: postId,
            count: 10
        } as types.PostsGetLikesConnection).then((resp: types.PostsLikesConnection) => {
            for (let user of resp.items) {
                store.dispatch({type: "users/set", user: user} as SetUser)
            }

            const likers: string[] = resp.items.map(u => u.id);
            store.dispatch({
                type: "posts/appendLikers",
                users: likers,
                postId: postId
            } as AppendLikers);

            resolve();
        })
    })
}
