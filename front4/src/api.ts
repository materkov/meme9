import {Query, QueryParams} from "./types";
import {writeStore} from "./store";

export function api(query: QueryParams): Promise<Query> {
    return new Promise((resolve, reject) => {
        let origin = window.location.origin;
        if (origin == "http://localhost:3000") {
            origin = "http://localhost:8000";
        }

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