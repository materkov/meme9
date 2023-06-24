import * as types from "../types/types";
import {create} from "zustand";
import {articlesLastPosted} from "./api";

export interface DiscoverPage {
    articles: types.Article[]
    fetch: () => void
}

export const useDiscoverPage = create<DiscoverPage>()(set => ({
    articles: [],
    fetch: () => {
        articlesLastPosted().then(articles =>
            set({articles: articles})
        );
    }
}));
