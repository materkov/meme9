import * as types from "../types/types";

function api<T>(method: string, args: any): Promise<T> {
    return new Promise((resolve, reject) => {
        fetch("/api/" + method, {
            method: 'POST',
            body: JSON.stringify(args),
            headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            }
        })
            .then(r => r.json())
            .then(r => {
                if (r.error) {
                    reject();
                } else {
                    resolve(r)
                }
            })
            .catch(reject)
    })
}

export function articlesSave(req: types.ArticlesSave): Promise<void> {
    return api("articles.save", req)
}

export function postsAdd(req: types.PostsAddReq): Promise<void> {
    return api("posts.add", req)
}

export function postsList(): Promise<types.Post[]> {
    return api("posts.list", {})
}

export function articlesList(req: types.ArticlesList): Promise<types.Article> {
    return api("articles.list", req)
}

export function articlesLastPosted(): Promise<types.Article[]> {
    return api('articles.lastPosted', {})
}

export function articlesListPostedByUser(req: { userId: string }): Promise<types.Article[]> {
    return api("articles.listPostedByUser", req);
}

export function usersList(req: { userIds: string[] }): Promise<types.User[]> {
    return api("users.list", req);
}
