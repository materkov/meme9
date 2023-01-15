import {api2, UsersFollow, UsersUnfollow} from "../../store/types";
import {store} from "../store";
import {SetIsFollowing} from "../reducers";

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
