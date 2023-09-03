import {create} from "zustand";
import {Post, postsListPostedByUser, User, usersList} from "../api/api";

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

        postsListPostedByUser({"userId": userId}).then(r => {
            set({posts: r});
        })
        usersList({"userIds": [userId]}).then(r => {
            set({user: r[0]})
        })
    },
}));
