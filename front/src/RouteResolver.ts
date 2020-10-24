import * as schema from "./schema/login";

let cachedRoutes: { [url: string]: schema.ResolveRouteResponse } = {};
let waiting: { [url: string]: ((r: schema.ResolveRouteResponse) => void)[] } = {};

if (window.InitData) {
    cachedRoutes[window.location.pathname] = {
        request: window.InitRequest,
        js: window.InitJsBundles,
    };
}

export function resolveRoute(url: string): Promise<schema.ResolveRouteResponse> {
    return new Promise<schema.ResolveRouteResponse>((resolve, reject) => {
        let path = url;

        if (!path.startsWith("/")) {
            const urlObject = new URL(url);
            if (!urlObject || !urlObject.pathname) {
                reject();
                return;
            }

            path = urlObject.pathname;
        }

        const cachedData = cachedRoutes[path];
        if (cachedData) {
            resolve(cachedData);
            return;
        }

        const waitingList = waiting[path];
        if (waitingList) {
            waitingList.push(resolve);
            return;
        }

        waiting[path] = [];

        fetch("/resolve-route", {
            method: 'POST',
            body: JSON.stringify({
                url: path,
            })
        }).then(r => r.json()).then((r: schema.ResolveRouteResponse) => {
            if (!r.request) {
                reject();
                return;
            }

            cachedRoutes[path] = r;
            resolve(r);

            const waitingList = waiting[path];
            if (waitingList) {
                for (let resolver of waitingList) {
                    resolver(r);
                }
                delete waiting[path];
            }
        }).catch(() => {
            reject();
        })
    })
}

