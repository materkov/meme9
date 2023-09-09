import {create} from "zustand";
import {Post, postsListPostedByUser, User, usersList} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";

export interface Profile {
    user: User;
    posts: Post[];
    fetched: { [id: string]: boolean };
    fetch: (userId: string) => void;
}

export const useProfile = create<Profile>()((set, get) => ({
    user: new User(),
    fetched: {},
    posts: [],
    fetch: (userId: string) => {
        if (get().fetched[userId]) {
            return
        }
        set({
            fetched: {
                ...get().fetched, [userId]: true,
            }
        });

        const prefetch = tryGetPrefetch('__userPage');
        if (prefetch && prefetch.user_id === userId) {
            set({
                posts: prefetch.posts,
                user: prefetch.user,
            })
            return;
        }

        postsListPostedByUser({"userId": userId}).then(r => {
            set({posts: r});
        })
        usersList({"userIds": [userId]}).then(r => {
            set({user: r[0]})
        })
    },
}));
