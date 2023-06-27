import * as types from "../types/types";
import {create} from "zustand";
import {articlesList} from "./api";

export interface ArticlePage {
    articles: { [id: string]: types.Article }
    error: any
    setText: (articleId: string, paragraphId: string, text: string) => void
    fetch: (id: string) => void
}

export const useArticlePage = create<ArticlePage>()((set, get) => ({
    articles: {},
    error: null,
    setText: (articleId: string, paragraphId: string, text: string) => set((state: ArticlePage) => {
        const copyArticle = structuredClone(state.articles[articleId]);
        const p = copyArticle.paragraphs.filter((p: types.Paragraph) => (p.text?.id == paragraphId));
        if (p[0].text) {
            p[0].text.text = text
        }

        return {...state.articles, [articleId]: copyArticle}
    }),
    fetch: (id: string) => {
        if (get().articles[id]) {
            return;
        }

        if ((window as any).__prefetchApi) {
            set({
                articles: {
                    ...get().articles, [(window as any).__prefetchApi.id]: (window as any).__prefetchApi
                }
            });
            delete (window as any).__prefetchApi;
            return;
        }

        articlesList({id: id})
            .then(data => {
                set({articles: {...get().articles, [data.id]: data}})
            })
            .catch(e => set({
                error: e
            }))
    },
}))
