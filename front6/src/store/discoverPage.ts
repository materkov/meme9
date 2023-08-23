import * as types from "../types/types";
import {create} from "zustand";
import {postsList} from "./api";

export interface DiscoverPage {
    posts: types.Post[]
    fetched: boolean
    fetch: () => void
    refetch: () => void
}

export const useDiscoverPage = create<DiscoverPage>()((set, get) => ({
    posts: [],
    fetched: false,
    fetch: () => {
        if (get().fetched) return;

        set(() => ({
            fetched: true,
        }));

        get().refetch();
    },
    refetch: () => {
        postsList().then(posts =>
            set({posts: posts})
        );
    }
}));
