import {create} from "zustand";
import {Post, postsListPostedByUser, usersList} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";
import {useResources} from "./resources";

export interface Profile {
    postIds: { [id: string]: string[] };
    fetched: { [id: string]: boolean };
    cursors: { [id: string]: string };
    fetch: (userId: string) => void;
    fetchMore: (userId: string) => void;
}

export const useProfile = create<Profile>()((set, get) => ({
    fetched: {},
    postIds: {},
    cursors: {},
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
            prefetch.posts.items.map((post: Post) => useResources.getState().setPost(post));

            set({
                postIds: {
                    ...get().postIds,
                    [prefetch.user_id]: prefetch.posts.items.map((post: Post) => post.id)
                },
                cursors: {
                    ...get().cursors,
                    [prefetch.user_id]: prefetch.posts.pageToken || "",
                }
            })
            return;
        }

        postsListPostedByUser({"userId": userId, after: ""}).then(r => {
            set({
                postIds: {
                    ...get().postIds,
                    [userId]: r.items.map((post: Post) => post.id)
                },
                cursors: {
                    ...get().cursors,
                    [userId]: r.pageToken || "",
                }
            })
            r.items.map((post: Post) => useResources.getState().setPost(post));
        })
        usersList({"userIds": [userId]}).then(r => {
            r.map(user => useResources.getState().setUser(user));
        })
    },
    fetchMore: userId => {
        const after = get().cursors[userId];
        if (!after) {
            return;
        }

        postsListPostedByUser({"userId": userId, after}).then(r => {
            set({
                postIds: {
                    ...get().postIds,
                    [userId]: [
                        ...get().postIds[userId],
                        ...r.items.map((post: Post) => post.id),
                    ],
                },
                cursors: {
                    ...get().cursors,
                    [userId]: r.pageToken || "",
                }
            })
            r.items.map((post: Post) => useResources.getState().setPost(post));
        })
    }
}));
