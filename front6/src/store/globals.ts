import {create} from "zustand";

export interface Globals {
    authToken: string;
    viewerId: string;
    viewerName: string;
}

export const useGlobals = create<Globals>()(set => {
    let authToken = window.__prefetchApi.authToken;
    let viewerId = window.__prefetchApi.viewerId;
    let viewerName = window.__prefetchApi.viewerName;
    if (authToken) {
        // TODO think about it
        localStorage.setItem('authToken', authToken);
    } else {
        localStorage.removeItem('authToken');
    }

    return {
        authToken: authToken,
        viewerId: viewerId,
        viewerName: viewerName,
    }
});
