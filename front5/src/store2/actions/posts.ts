import {api2, Post, PostLikeData, PostsAdd, PostsDelete, PostsLike, PostsUnlike} from "../../store/types";
import {Global, store} from "../store";
import {AppendFeed, DeleteFromFeed, SetLikes, SetPost} from "../reducers";

export function deletePost(postId: string) {
    api2("posts.delete", {
        id: postId,
    } as PostsDelete);

    store.dispatch({type: "feed/delete", postId: postId} as DeleteFromFeed)
}

export function add(data: PostsAdd): Promise<void> {
    return new Promise((resolve, reject) => {
        api2("posts.add", data).then((post: Post) => {
            store.dispatch({type: "posts/set", post: post} as SetPost)
            store.dispatch({type: "feed/append", items: [post.id], prepend: true} as AppendFeed)

            return resolve();
        }).catch(() => {
            return resolve();
        })
    })
}

export function like(postId: string) {
    const state = store.getState() as Global;

    store.dispatch({
        type: 'posts/setLikes',
        postId: postId,
        count: (state.posts.likesCount[postId] || 0) + 1,
        isLiked: true,
    } as SetLikes);

    api2("posts.like", {postId: postId} as PostsLike).then((resp: PostLikeData) => {
        store.dispatch({
            type: 'posts/setLikes',
            postId: postId,
            count: resp.totalCount || 0,
            isLiked: true,
        } as SetLikes);
    });
}

export function unlike(postId: string) {
    const state = store.getState() as Global;

    store.dispatch({
        type: 'posts/setLikes',
        postId: postId,
        count: (state.posts.likesCount[postId] || 0) - 1,
        isLiked: false,
    } as SetLikes);

    api2("posts.unlike", {postId: postId} as PostsUnlike).then((resp: PostLikeData) => {
        store.dispatch({
            type: 'posts/setLikes',
            postId: postId,
            count: resp.totalCount || 0,
            isLiked: false,
        } as SetLikes);
    });
}
