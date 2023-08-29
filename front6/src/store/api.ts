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

export function postsAdd(req: types.PostsAddReq): Promise<void> {
    return api("posts.add", req)
}

export function postsList(): Promise<types.Post[]> {
    return api("posts.list", {})
}

export function postsListPostedByUser(req: types.PostsListPostedByUser): Promise<types.Post[]> {
    return api("posts.listPostedByUser", req)
}

export function postsListById(req: types.PostsListById): Promise<types.Post | undefined> {
    return api("posts.listById", req)
}

export function usersList(req: { userIds: string[] }): Promise<types.User[]> {
    return api("users.list", req);
}
