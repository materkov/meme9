export type User = {
    id: string;
    name?: string;
    posts?: Post[];
}

export type Post = {
    id: string;
    userId: string;
    user?: User;
    date: string;
    text: string;
}

export function api(url: string, params: any = {}): Promise<any> {
    const body = new FormData();

    for (let key in params) {
        body.append(key, params[key]);
    }

    const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me/api" : "http://localhost:8000/api";

    return new Promise((resolve, reject) => {
        return fetch(apiHost + url, {
            method: 'POST',
            body: body,
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken'),
            },
        })
            .then(r => r.json())
            .then(r => {
                if (r.ok) {
                    resolve(r.data);
                } else {
                    reject(r.error);
                }
            })
    })
}