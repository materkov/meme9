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

//export let apiHost = "http://localhost:8000/api";
export let apiHost = "https://meme.mmaks.me/api";

export function api(url: string, params: any = {}): Promise<any> {
    const body = new FormData();

    for (let key in params) {
        body.append(key, params[key]);
    }

    return fetch(apiHost + url, {
        method: 'POST',
        body: body,
        headers: {
            'authorization': 'Bearer ' + localStorage.getItem('authToken'),
        },
    })
        .then(r => r.json());
}