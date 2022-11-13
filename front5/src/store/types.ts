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
    photoUrl?: string;
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
