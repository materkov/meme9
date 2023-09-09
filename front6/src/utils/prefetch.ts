export function tryGetPrefetch(key: string) {
    if (!window.__prefetchApi[key]) {
        return undefined;
    }

    const data = window.__prefetchApi[key];
    delete window.__prefetchApi[key];

    return data;
}
