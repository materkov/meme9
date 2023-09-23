import {create} from "zustand";
import {Post, postsListPostedByUser, User, usersList} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";

export interface Profile {
    users: { [id: string]: User };
    posts: { [id: string]: Post[] };
    fetched: { [id: string]: boolean };
    fetch: (userId: string) => void;
    setStatus: (userId: string, status: string) => void;
}

export const useProfile = create<Profile>()((set, get) => ({
    users: {},
    fetched: {},
    posts: {},
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
                posts: {
                    [prefetch.user_id]: prefetch.posts,
                },
                users: {
                    [prefetch.user_id]: prefetch.user,
                },
            })
            return;
        }

        postsListPostedByUser({"userId": userId}).then(r => {
            set({
                posts: {
                    ...get().posts,
                    [userId]: r
                }
            });
        })
        usersList({"userIds": [userId]}).then(r => {
            set({
                users: {
                    ...get().users,
                    [userId]: r[0]
                },
            });
        })
    },
    setStatus: (userId, status) => {
        if (!get().users[userId]) {
            return;
        }

        set({
            users: {
                ...get().users,
                [userId]: {
                    ...get().users[userId],
                    status: status
                },
            }
        });
    }
}));
