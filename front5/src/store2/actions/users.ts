import {
    api2,
    PostsList,
    User,
    UsersEdit,
    UsersFollow,
    UsersPostsList,
    UsersSetAvatar,
    UsersUnfollow
} from "../../store/types";
import {Global, store} from "../store";
import {parsePostsList} from "../helpers/posts";
import {AppendPosts, SetIsFollowing, SetPostsCount, SetUser} from "../reducers/users";

export function follow(userId: string) {
    store.dispatch({
        type: 'users/setIsFollowing',
        userId: userId,
        isFollowing: true,
    } as SetIsFollowing);

    api2("users.follow", {userId: userId} as UsersFollow);
}

export function unfollow(userId: string) {
    store.dispatch({
        type: 'users/setIsFollowing',
        userId: userId,
        isFollowing: false,
    } as SetIsFollowing);

    api2("users.unfollow", {userId: userId} as UsersUnfollow);
}

export function edit(data: UsersEdit): Promise<void> {
    return new Promise((resolve, reject) => {
        const state = store.getState() as Global;

        const user = {...state.users.byId[data.userId]};
        user.name = data.name;

        store.dispatch({
            type: 'users/set',
            user: user,
        } as SetUser);

        api2("users.edit", data).then(() => resolve());
    })
}

export function usersSetOnline() {
    api2("users.setOnline", {});
}

export function usersSetAvatar(uploadToken: string): Promise<void> {
    return new Promise((resolve, reject) => {
        api2("users.setAvatar", {uploadToken: uploadToken} as UsersSetAvatar)
            .then((user: User) => {
                store.dispatch({
                    type: 'users/set',
                    user: user,
                } as SetUser);
                resolve();
            })
    })
}

export function loadUserPage(userId: string): Promise<undefined> {
    return new Promise((resolve, reject) => {
        api2("users.posts.list", {userId: userId, count: 10} as UsersPostsList).then((resp: PostsList) => {
            parsePostsList(resp);

            store.dispatch({
                type: "users/appendPosts",
                userId: userId,
                posts: resp.items.map(post => post.id)
            } as AppendPosts);

            store.dispatch({
                type: "users/setPostsCount",
                userId: userId,
                count: resp.totalCount || 0,
            } as SetPostsCount);

            resolve(undefined);
        })
    })
}

export function loadUserPostsCount(userId: string): Promise<void> {
    return new Promise((resolve, reject) => {
        const state = store.getState() as Global;
        if (state.users.postsCount.hasOwnProperty(userId)) {
            resolve();
            return
        }

        store.dispatch({
            type: "users/setPostsCount",
            userId: userId,
            count: null,
        } as SetPostsCount);

        api2("users.posts.list", {userId: userId, count: 0} as UsersPostsList).then((resp: PostsList) => {
            store.dispatch({
                type: "users/setPostsCount",
                userId: userId,
                count: resp.totalCount || 0,
            } as SetPostsCount);

            resolve();
        })
    })
}
