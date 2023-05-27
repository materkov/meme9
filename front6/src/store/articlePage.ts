import * as types from "../types/types";
import {create} from "zustand";
import {api} from "./api";

export interface ArticlePage {
    article: types.Article
    error: any
    setArticle: (article: types.Article) => void
    setText: (paragraphId: string, text: string) => void
    fetch: (id: string) => void
}

export const useArticlePage = create<ArticlePage>()(set => ({
    article: new types.Article(),
    error: null,
    setArticle: (article: types.Article) => set(() => ({article: article})),
    setText: (paragraphId: string, text: string) => set((state: ArticlePage) => {
        const copyArticle = structuredClone(state.article);
        const p = copyArticle.paragraphs.filter((p: types.Paragraph) => (p.text?.id == paragraphId));
        if (p[0].text) {
            p[0].text.text = text
        }

        return {article: copyArticle}
    }),
    fetch: (id: string) => {
        if (window.__prefetchApi) {
            set({article: window.__prefetchApi});
            delete window.__prefetchApi;
            return;
        }

        api<types.Article>('articles.list', {id: id})
            .then(data => set({
                article: data,
            }))
            .catch(e => set({
                error: e
            }))
    },
}))
