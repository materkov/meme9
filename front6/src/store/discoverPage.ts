import {create} from "zustand";
import {FeedType, Post, postsList} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";
import {useResources} from "./resources";

export interface DiscoverPage {
    posts: string[]
    fetched: boolean
    fetch: () => void
    refetch: () => void
    type: FeedType
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
            set({posts: prefetch.map((post: Post) => post.id)});
            prefetch.map((post: Post) => useResources.getState().setPost(post));
            return;
        }

        get().refetch();
    },
    refetch: () => {
        postsList({type: get().type}).then(posts => {
            set({posts: posts.map(post => post.id)});
            posts.map(post => useResources.getState().setPost(post));
        });
    }
}));
