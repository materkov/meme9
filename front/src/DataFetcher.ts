import * as schema from "./schema/login";

let cache: { [request: string]: schema.AnyRenderer } = {};

if (window.InitData) {
    cache[JSON.stringify(window.InitRequest)] = window.InitData;
}

export function fetchData(request: schema.AnyRequest): Promise<schema.AnyRenderer> {
    return new Promise<schema.AnyRenderer>((resolve, reject) => {
        const requestText = JSON.stringify(request);

        const cached = cache[requestText];
        if (cached) {
            resolve(cached);
            return;
        }

        fetch('/api', {
            method: 'POST',
            body: requestText,
        }).then(r => (
            r.json()
        )).then(r => {
            cache[requestText] = r;
            resolve(r);
        }).catch(() => {
            reject();
        });
    });
}
