import {create} from "zustand";
import {postsListById} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";
import {useResources} from "./resources";

export interface PostPage {
    errors: { [id: string]: any }
    fetch: (id: string) => void
}

export const usePostPage = create<PostPage>()((set, get) => ({
    errors: {},
    fetch: (id: string) => {
        if (useResources.getState().posts[id]) {
            return;
        }

        const prefetch = tryGetPrefetch('__postPagePost');
        if (prefetch) {
            useResources.getState().setPost(prefetch);
            return;
        }

        postsListById({id: id})
            .then(data => {
                useResources.getState().setPost(data);
            })
            .catch(e => set({
                errors: {...get().errors, [id]: e}
            }))
    },
}))
