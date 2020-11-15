import * as schema from "./schema/api";

let cachedRoutes: { [url: string]: schema.ResolveRouteResponse } = {};
let waiting: { [url: string]: ((r: schema.ResolveRouteResponse) => void)[] } = {};

if (window.InitApiResponse) {
    cachedRoutes[window.location.pathname] = {
        apiMethod: window.InitApiMethod,
        apiRequest: JSON.stringify(window.InitApiRequest),
        js: window.InitJsBundles,
        rootComponent: window.InitRootComponent,
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

