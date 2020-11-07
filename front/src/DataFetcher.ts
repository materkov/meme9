import * as schema from "./schema/login";

let cache: { [request: string]: any } = {};

if (window.InitApiMethod) {
    const cacheKey = window.InitApiMethod + "__" + JSON.stringify(window.InitApiRequest);
    cache[cacheKey] = window.InitApiResponse;
}

export function fetchData(method: string, args: string): Promise<any> {
    return new Promise<any>((resolve, reject) => {
        //const argsText = JSON.stringify(args);
        const argsText = args;

        const cacheKey = method + "__" + argsText;
        const cached = cache[cacheKey];
        if (cached) {
            resolve(cached);
            return;
        }

        fetch('/api/' + method, {
            method: 'POST',
            headers: {
                'x-csrf-token': window.CSRFToken,
            },
            body: argsText,
        }).then(r => (
            r.json()
        )).then(r => {
            cache[cacheKey] = r;
            resolve(r);
        }).catch(() => {
            reject();
        });
    });
}
