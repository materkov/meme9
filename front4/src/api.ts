import {Query, QueryParams} from "./types";
import {writeStore} from "./store/store";

function getApiOrigin(): string {
    let origin = window.location.origin;
    if (origin == "http://localhost:3000") {
        origin = "http://localhost:8000";
    }

    return origin;
}

export function api(query: QueryParams): Promise<Query> {
    return new Promise((resolve, reject) => {
        let origin = getApiOrigin();

        let headers: Record<string, string> = {
            'Content-Type': 'application/json',
        }

        const token = localStorage.getItem('authToken');
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }

        fetch(origin + "/gql", {
            method: 'POST',
            headers: headers,
            body: JSON.stringify(query),
        })
            .then(data => data.json())
            .then(data => {
                resolve(data as Query);
                writeStore(data);
            })
            .catch(() => {
                reject();
            })
    })
}

export function apiUpload(file: File) {
    return new Promise((resolve, reject) => {
        let origin = getApiOrigin();

        let headers: Record<string, string> = {
            'Content-Type': 'application/octet-stream',
        };

        const token = localStorage.getItem('authToken');
        if (token) {
            headers['Authorization'] = 'Bearer ' + token;
        }

        fetch(origin + "/upload", {
            method: 'POST',
            headers: headers,
            body: file,
        })
            .then(data => data.json())
            .then(data => {
                resolve(data);
            })
            .catch(() => {
                reject();
            })
    })
}
