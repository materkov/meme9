import * as types from "../types/types";
import {create} from "zustand";
import {api} from "./api";

export interface DiscoverPage {
    articles: types.Article[]
    fetch: () => void
}

export const useDiscoverPage = create<DiscoverPage>()(set => ({
    articles: [],
    fetch: () => {
        api<types.Article[]>('articles.lastPosted', {}).then(articles =>
            set({articles: articles})
        );
    }
}));
