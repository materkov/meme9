import {create} from "zustand";
import {Post, postsListById} from "../api/api";

export interface PostPage {
    posts: { [id: string]: Post | undefined }
    error: any
    fetch: (id: string) => void
}

export const usePostPage = create<PostPage>()((set, get) => ({
    posts: {},
    error: null,
    fetch: (id: string) => {
        if (get().posts[id]) {
            return;
        }

        if (window.__prefetchApi.__postPagePost) {
            set({
                posts: {
                    ...get().posts, [window.__prefetchApi.__postPagePost.id]: (window as any).__prefetchApi.__postPagePost
                }
            });
            delete window.__prefetchApi.__postPagePost;
            return;
        }

        postsListById({id: id})
            .then(data => {
                set({posts: {...get().posts, [id]: data}})
            })
            .catch(e => set({
                error: e
            }))
    },
}))
