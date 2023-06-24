import * as types from "../types/types";
import {create} from "zustand";
import {articlesListPostedByUser, usersList} from "./api";

export interface Profile {
    user: types.User;
    articles: types.Article[];
    setArticles: (articles: types.Article[]) => void;
    setUser: (user: types.User) => void;
    fetch: (userId: string) => void;
}

export const useProfile = create<Profile>()(set => ({
    user: new types.User(),
    articles: [],
    setArticles: (articles: types.Article[]) => set(() => ({
        articles: articles
    })),
    setUser: (user: types.User) => set(() => ({
        user: user,
    })),
    fetch: (userId: string) => {
        articlesListPostedByUser({"userId": userId}).then(r => {
            set({articles: r});
        })
        usersList({"userIds": [userId]}).then(r => {
            set({user: r[0]})
        })
    },
}));
