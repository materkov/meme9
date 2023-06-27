import * as types from "../types/types";
import {create} from "zustand";
import {articlesListPostedByUser, usersList} from "./api";

export interface Profile {
    user: types.User;
    articles: types.Article[];
    fetched: { [id: string]: boolean };
    fetch: (userId: string) => void;
}

export const useProfile = create<Profile>()((set, get) => ({
    user: new types.User(),
    fetched: {},
    articles: [],
    fetch: (userId: string) => {
        if (get().fetched[userId]) {
            return
        }
        set({
            fetched: {
                ...get().fetched, [userId]: true,
            }
        });

        articlesListPostedByUser({"userId": userId}).then(r => {
            set({articles: r});
        })
        usersList({"userIds": [userId]}).then(r => {
            set({user: r[0]})
        })
    },
}));
