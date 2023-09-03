function api<T>(method: string, args: any): Promise<T> {
    let headers: Record<string, string> = {};
    if (window.__prefetchApi.authToken) {
        headers['authorization'] = 'Bearer ' + window.__prefetchApi.authToken;
    }

    return new Promise((resolve, reject) => {
        fetch('/api/' + method, {
            method: 'POST',
            body: JSON.stringify(args),
            headers: headers,
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

export class User {
    id = ""
    name = ""
}

export class Post {
    id: string = ""
    userId: string = ""
    date: string = ""
    text: string = ""
    user?: User = undefined
}

export class Void {
}

export class PostsAddReq {
    text: string = ""
}

export function postsAdd(req: PostsAddReq): Promise<void> {
    return api("posts.add", req)
}

export function postsList(): Promise<Post[]> {
    return api("posts.list", {})
}

export class PostsListPostedByUserReq {
    userId: string = ""
}

export function postsListPostedByUser(req: PostsListPostedByUserReq): Promise<Post[]> {
    return api("posts.listPostedByUser", req)
}

export class PostsListByIdReq {
    id: string = ""
}

export function postsListById(req: PostsListByIdReq): Promise<Post | undefined> {
    return api("posts.listById", req)
}

export class PostsDeleteReq {
    postId: string = ""
}

export function postsDelete(req: PostsDeleteReq): Promise<Void> {
    return api("posts.delete", req);
}

export class UsersListReq {
    userIds: string[] = []
}

export function usersList(req: UsersListReq): Promise<User[]> {
    return api("users.list", req);
}
