import * as types from "../types/types";
import {create} from "zustand";
import {articlesLastPosted} from "./api";

export interface DiscoverPage {
    articles: types.Article[]
    fetched: boolean
    fetch: () => void
}

export const useDiscoverPage = create<DiscoverPage>()((set, get) => ({
    articles: [],
    fetched: false,
    fetch: () => {
        if (get().fetched) return;

        set(() => ({
            fetched: true,
        }));

        articlesLastPosted().then(articles =>
            set({articles: articles})
        );
    }
}));
