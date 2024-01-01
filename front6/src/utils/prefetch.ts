export function tryGetPrefetch(key: string) {
    if (!window.__prefetchApi[key]) {
        return undefined;
    }

    const data = window.__prefetchApi[key];
    delete window.__prefetchApi[key];

    return data;
}

// TODO fix prefetch

export function usePrefetch(key: string, prefetch: (data: any) => void) {
    if (window.__prefetchApi[key]) {

        prefetch(window.__prefetchApi[key]);
        delete window.__prefetchApi[key];
    }
}
