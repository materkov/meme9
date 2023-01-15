export type User = {
    id: string;
    name?: string;
    avatar?: string;
    bio?: string;
}

export type Post = {
    id: string;
    userId: string;
    date: string;
    text: string;
    canDelete?: boolean;
    isDeleted?: boolean;
    photoId?: string;
}

export type Photo = {
    id: string;
    url: string;
    address: string;
    width: number;
    height: number;
    thumbs: PhotoThumb[];
}

export type PhotoThumb = {
    width: number;
    height: number;
    address: string;
}

export function api(url: string, params: any = {}): Promise<any> {
    const body = new FormData();
    body.set("token", localStorage.getItem('authToken') || '');

    for (let key in params) {
        body.append(key, params[key]);
    }

    const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me/api" : "http://localhost:8000/api";

    return new Promise((resolve, reject) => {
        return fetch(apiHost + url, {
            method: 'POST',
            body: body,
        })
            .then(resp => {
                if (resp.status == 200) {
                    resp.json().then(resolve);
                } else {
                    resp.text().then(reject);
                }
            })
    })
}

export function api2(method: string, args: any): Promise<any> {
    const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me/api2/" : "http://localhost:8000/api2/";

    return new Promise((resolve, reject) => {
        return fetch(apiHost + method , {
            method: 'POST',
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken') || '',
            },
            body: JSON.stringify(args),
        })
            .then(resp => resp.json())
            .then(resp => {
                if (resp.error) {
                    reject(resp.error);
                } else {
                    resolve(resp.data);
                }
            })
    })
}

export function uploadApi(file: File): Promise<string> {
    const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me" : "http://localhost:8000";

    return new Promise((resolve, reject) => {
        return fetch(apiHost + "/upload", {
            method: 'POST',
            headers: {
                'content-type': 'application/octet-stream',
            },
            body: file,
        })
            .then(resp => {
                if (resp.status == 200) {
                    resp.json().then((res) => resolve(res.uploadToken));
                } else {
                    resp.text().then(reject);
                }
            })
    })
}

export interface PostLikeData {
    isViewerLiked: boolean
    totalCount: number
    items: string
}

export interface Edges {
    totalCount: number
    items: string[]
    nextCursor: string
}

export interface Viewer {
    url: string;
    viewerId: string;
}

export interface Online {
    url: string;
    isOnline: boolean;
}

export interface PostsAdd {
    text: string;
    photo: string;
}

export interface PostsDelete {
    id: string;
}

export interface PostsLike {
    postId: string;
}

export interface PostsUnlike {
    postId: string;
}
