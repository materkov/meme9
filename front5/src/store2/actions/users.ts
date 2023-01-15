import {api2, UsersEdit, UsersFollow, UsersUnfollow} from "../../store/types";
import {Global, store} from "../store";
import {SetIsFollowing, SetUser} from "../reducers";

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
