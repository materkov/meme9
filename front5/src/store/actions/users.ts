import * as types from "../../api/types";
import {LoadingState, store} from "../store";
import {parsePostsList} from "../helpers/posts";
import {api} from "../../api/api";

export function follow(userId: string) {
    store.dispatch({
        type: 'users/setIsFollowing',
        userId: userId,
        isFollowing: true,
    });

    const oldCount = store.getState().users.followersCount[userId] || 0;
    store.dispatch({
        type: 'users/setFollowersCount',
        userId: userId,
        count: oldCount + 1,
        isViewerFollowing: true,
    });

    api("users.follow", {userId: userId} as types.UsersFollow);
}

export function unfollow(userId: string) {
    store.dispatch({
        type: 'users/setIsFollowing',
        userId: userId,
        isFollowing: false,
    });

    const oldCount = store.getState().users.followersCount[userId] || 0;
    store.dispatch({
        type: 'users/setFollowersCount',
        userId: userId,
        count: oldCount - 1,
        isViewerFollowing: false,
    });

    api("users.unfollow", {userId: userId} as types.UsersUnfollow);
}

export function edit(data: types.UsersEdit): Promise<void> {
    return new Promise((resolve, reject) => {
        const state = store.getState();

        const user = {...state.users.byId[data.userId]};
        user.name = data.name;

        store.dispatch({
            type: 'users/set',
            user: user,
        });

        api("users.edit", data).then(() => resolve());
    })
}

export function usersSetOnline() {
    api("users.setOnline", {});
}

export function usersSetAvatar(uploadToken: string): Promise<void> {
    return new Promise((resolve, reject) => {
        api("users.setAvatar", {uploadToken: uploadToken} as types.UsersSetAvatar)
            .then((user: types.User) => {
                store.dispatch({
                    type: 'users/set',
                    user: user,
                });
                resolve();
            })
    })
}

export function fetchUserPosts(userId: string): Promise<void> {
    return new Promise((resolve) => {
        const key = "fetchUserPosts:" + userId;
        if (store.getState().routing.fetchLockers[key]) {
            return;
        }
        store.dispatch({type: "routes/setFetchLocker", key: key, state: LoadingState.LOADING});

        api("users.posts.list", {userId: userId, count: 10} as types.UsersPostsList).then((resp: types.PostsList) => {
            parsePostsList(resp);

            store.dispatch({
                type: "users/appendPosts",
                userId: userId,
                posts: resp.items.map(post => post.id)
            });

            store.dispatch({
                type: "users/setPostsCount",
                userId: userId,
                count: resp.totalCount || 0,
            });

            resolve();
        })
    })
}

export function loadUserPostsCount(userId: string): Promise<void> {
    return new Promise((resolve, reject) => {
        const state = store.getState();
        if (state.users.postsCount[userId] !== undefined) {
            resolve();
            return
        }

        store.dispatch({
            type: "users/setPostsCount",
            userId: userId,
            count: null,
        });

        api("users.posts.list", {userId: userId, count: 0} as types.UsersPostsList).then((resp: types.PostsList) => {
            store.dispatch({
                type: "users/setPostsCount",
                userId: userId,
                count: resp.totalCount || 0,
            });

            resolve();
        })
    })
}

export function fetchFollowersCount(userId: string) {
    const key = "fetchFollowersCount:" + userId;

    if (store.getState().routing.fetchLockers[key]) {
        return;
    }
    store.dispatch({type: "routes/setFetchLocker", key: key, state: LoadingState.LOADING});

    api("users.followers.list", {
        userId: userId,
        count: 0
    } as types.UserFollowersCount).then((resp: types.FollowersEdges) => {
        store.dispatch({
            type: "users/setFollowersCount",
            userId: userId,
            count: resp.totalCount || 0,
            isViewerFollowing: resp.isFollowing || false,
        });
    })
}

export function fetchFollowingCount(userId: string) {
    const key = "fetchFollowingCount:" + userId;

    if (store.getState().routing.fetchLockers[key]) {
        return;
    }
    store.dispatch({type: "routes/setFetchLocker", key: key, state: LoadingState.LOADING});

    api("users.following.list", {
        userId: userId,
        count: 0
    } as types.UserFollowersCount).then((resp: types.Edges) => {
        store.dispatch({
            type: "users/setFollowingCount",
            userId: userId,
            count: resp.totalCount || 0,
        });
    })
}
