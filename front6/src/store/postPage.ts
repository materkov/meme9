import {create} from "zustand";
import {Post, postsListById} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";

export interface PostPage {
    posts: { [id: string]: Post | undefined }
    errors: { [id: string]: any }
    fetch: (id: string) => void
}

export const usePostPage = create<PostPage>()((set, get) => ({
    posts: {},
    errors: {},
    fetch: (id: string) => {
        if (get().posts[id]) {
            return;
        }

        const prefetch = tryGetPrefetch('__postPagePost');
        if (prefetch) {
            set({
                posts: {
                    ...get().posts, [prefetch.id]: prefetch,
                }
            });
            return;
        }

        postsListById({id: id})
            .then(data => {
                set({posts: {...get().posts, [id]: data}})
            })
            .catch(e => set({
                errors: {...get().errors, [id]: e}
            }))
    },
}))
