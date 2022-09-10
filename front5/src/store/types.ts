export type User = {
    id: string;
    name?: string;
    posts?: Post[];
    avatar?: string;
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
                    reject();
                }
            });
    })
}