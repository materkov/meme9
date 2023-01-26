import {store} from "../store/store";
import {ApiRequest} from "./types";

export function api(method: string, args: ApiRequest): Promise<any> {
    const apiHost = location.host == "meme.mmaks.me" ? "https://meme.mmaks.me/api2?method=" : "http://localhost:8000/api2?method=";
    const state = store.getState();

    return new Promise((resolve, reject) => {
        return fetch(apiHost + method, {
            method: 'POST',
            headers: {
                'authorization': 'Bearer ' + state.routing.accessToken,
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

    const form = new FormData();
    form.set("file", file);

    return new Promise((resolve, reject) => {
        return fetch(apiHost + "/upload", {
            method: 'POST',
            headers: {},
            body: form,
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
