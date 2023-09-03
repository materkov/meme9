import {create} from "zustand";
import {Post, postsList} from "../api/api";

export interface DiscoverPage {
    posts: Post[]
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

        if (window.__prefetchApi.__postsList) {
            set({posts: window.__prefetchApi.__postsList});
            delete window.__prefetchApi.__postsList;
            return;
        }

        get().refetch();
    },
    refetch: () => {
        postsList().then(posts =>
            set({posts: posts})
        );
    }
}));
