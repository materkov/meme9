import * as schema from "./schema/login";

let cache: { [request: string]: any } = {};

if (window.InitApiCommand) {
    const cacheKey = window.InitApiCommand + "__" + JSON.stringify(window.InitApiArgs);
    cache[cacheKey] = window.InitData;
}

export function fetchData(command: string, args: string): Promise<any> {
    return new Promise<any>((resolve, reject) => {
        //const argsText = JSON.stringify(args);
        const argsText = args;

        const cacheKey = command + "__" + argsText;
        const cached = cache[cacheKey];
        if (cached) {
            resolve(cached);
            return;
        }

        fetch('/api2/' + command, {
            method: 'POST',
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
