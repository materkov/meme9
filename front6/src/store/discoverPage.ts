import {create} from "zustand";
import {FeedType, Post, postsList} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";

export interface DiscoverPage {
    posts: Post[]
    type: FeedType
    fetched: boolean
    fetch: () => void
    refetch: () => void
    setType: (type: FeedType) => void
}

export const useDiscoverPage = create<DiscoverPage>()((set, get) => ({
    posts: [],
    fetched: false,
    type: FeedType.DISCOVER,
    setType: (type: FeedType) => {
        set({type});
    },
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
        postsList({type: get().type}).then(posts =>
            set({posts: posts})
        );
    }
}));
