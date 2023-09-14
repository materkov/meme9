import {cookieAuthToken, getCookie} from "../utils/cookie";

function api<T>(method: string, args: any): Promise<T> {
    let headers: Record<string, string> = {};

    let token = getCookie(cookieAuthToken);
    if (token) {
        headers['authorization'] = 'Bearer ' + token;
    }

    return new Promise((resolve, reject) => {
        // TODO think about this func
        fetch('/api/' + method, {
            credentials: 'omit',
            method: 'POST',
            body: JSON.stringify(args),
            headers: headers,
        })
            .then(r => {
                if (r.ok) {
                    r.json()
                        .then(resolve)
                        .catch(reject);
                } else {
                    r.json()
                        .then(r => {
                            reject(r.error);
                        })
                        .catch(reject);
                }
            })
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

export class AuthEmailReq {
    email: string = ""
    password: string = ""
}

export class AuthResp {
    token: string = ""
    userId: string = ""
    userName: string = ""
}

export function authLogin(req: AuthEmailReq): Promise<AuthResp> {
    return api("auth.login", req);
}

export function authRegister(req: AuthEmailReq): Promise<AuthResp> {
    return api("auth.register", req);
}

export class AuthVk {
    code: string = ""
    redirectUrl: string = ""
}

export function authVK(req: AuthVk): Promise<AuthResp> {
    return api("auth.vk", req);
}
