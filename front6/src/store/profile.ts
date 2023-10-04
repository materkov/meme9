import {create} from "zustand";
import {Post, postsListPostedByUser, usersList} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";
import {useResources} from "./resources";

export interface Profile {
    postIds: { [id: string]: string[] };
    fetched: { [id: string]: boolean };
    fetch: (userId: string) => void;
}

export const useProfile = create<Profile>()((set, get) => ({
    fetched: {},
    postIds: {},
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
            useResources.getState().setUser(prefetch.user);
            prefetch.posts.map((post: Post) => useResources.getState().setPost(post));

            set({
                postIds: {
                    ...get().postIds,
                    [prefetch.user_id]: prefetch.posts.map((post: Post) => post.id)
                }
            })
            return;
        }

        postsListPostedByUser({"userId": userId}).then(r => {
            set({
                postIds: {
                    ...get().postIds,
                    [userId]: r.map((post: Post) => post.id)
                }
            })
            r.map((post: Post) => useResources.getState().setPost(post));
        })
        usersList({"userIds": [userId]}).then(r => {
            r.map(user => useResources.getState().setUser(user));
        })
    },
}));
