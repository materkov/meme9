import {create} from "zustand";
import {FeedType, Post, PostsList, postsList} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";
import {useResources} from "./resources";

export interface DiscoverPage {
    posts: string[]
    postsPageToken: string
    fetched: boolean
    fetch: () => void
    refetch: () => void
    type: FeedType
    setType: (type: FeedType) => void
}

export const useDiscoverPage = create<DiscoverPage>()((set, get) => ({
    posts: [],
    postsPageToken: "",
    fetched: false,
    type: FeedType.DISCOVER,
    setType: (type: FeedType) => {
        set({type});
    },
    fetch: () => {
        // TODO think about refetching
        //if (get().fetched) return;

        set(() => ({
            fetched: true,
        }));

        postsList({
            count: 10,
            type: get().type,
            pageToken: get().postsPageToken,
        }).then(postsList => {
            set({
                posts: [
                    ...get().posts,
                    ...postsList.items.map(post => post.id),
                ],
                postsPageToken: postsList.pageToken,
            });
            postsList.items.map(post => useResources.getState().setPost(post));
        });
    },
    refetch: () => {
        const prefetch: PostsList = tryGetPrefetch('__postsList');
        if (prefetch) {
            set({
                posts: prefetch.items.map((post: Post) => post.id),
                postsPageToken: prefetch.pageToken,
            });
            prefetch.items.map((post: Post) => useResources.getState().setPost(post));
            return;
        }

        postsList({
            count: 10,
            type: get().type,
            pageToken: '',
        }).then(postsList => {
            set({
                posts: postsList.items.map(post => post.id),
                postsPageToken: postsList.pageToken,
            });
            postsList.items.map(post => useResources.getState().setPost(post));
        });
    },
}));
