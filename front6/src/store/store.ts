import * as types from "../types/types";
import {ArticlesSave} from "../types/types";
import {create} from "zustand";

const apiHost = document.location.host == "meme.mmaks.me" ?
    "https://meme.mmaks.me" : "http://localhost:8000";

export function getArticle(id: string): Promise<types.Article> {
    return new Promise((resolve, reject) => {
        fetch(apiHost + "/api/articles.list", {
            method: 'POST',
            body: JSON.stringify({"id": id})
        })
            .then(r => {
                r.json().then(data => {
                    if (r.status == 200) {
                        resolve(types.Article.fromJSON(data))
                    } else {
                        reject(data);
                    }
                }).catch(reject)
            }).catch(reject)
    })
}

export function saveArticle(article: ArticlesSave): Promise<undefined> {
    return new Promise((resolve, reject) => {
        fetch(apiHost + "/api/articles.save", {
            method: 'POST',
            body: JSON.stringify(article),
            headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            }
        })
            .then(r => r.json())
            .then(r => {
                if (r.error) {
                    reject();
                } else {
                    resolve(undefined)
                }
            })
            .catch(reject)
    })
}

export interface ArticlePage {
    article: types.Article
    setArticle: (article: types.Article) => void
    setText: (paragraphId: string, text: string) => void
}

export const useArticlePage = create<ArticlePage>()(set => ({
    article: new types.Article(),
    setArticle: (article: types.Article) => set(() => ({article: article})),
    setText: (paragraphId: string, text: string) => set((state: ArticlePage) => {
        const copyArticle = structuredClone(state.article);
        const p = copyArticle.paragraphs.filter((p: types.Paragraph) => (p.text?.id == paragraphId));
        if (p[0].text) {
            p[0].text.text = text
        }

        return {article: copyArticle}
    })
}))
