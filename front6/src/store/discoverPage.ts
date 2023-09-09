import {create} from "zustand";
import {Post, postsList} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";

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

        const prefetch = tryGetPrefetch('__postsList');
        if (prefetch) {
            set({posts: prefetch});
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
