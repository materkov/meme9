import {create} from "zustand";

export interface Globals {
    authToken: string;
    viewerId: string;
    viewerName: string;
    logout: () => void;
}

export const useGlobals = create<Globals>()(set => {
    let authToken = (window.__prefetchApi as any).authToken;
    if (authToken) {
        localStorage.setItem('authToken', authToken);
    } else {
        authToken = localStorage.getItem('authToken') || "";
    }

    let viewerId = (window.__prefetchApi as any).viewerId;
    if (viewerId) {
        localStorage.setItem('viewerId', viewerId);
    } else {
        viewerId = localStorage.getItem('viewerId') || "";
    }

    let viewerName = (window.__prefetchApi as any).viewerName;
    if (viewerName) {
        localStorage.setItem('viewerName', viewerName);
    } else {
        viewerName = localStorage.getItem('viewerName') || "";
    }

    return {
        authToken: authToken,
        viewerId: viewerId,
        viewerName: viewerName,
        logout: () => set(() => {
            localStorage.removeItem('viewerId');
            localStorage.removeItem('viewerName');

            return {
                viewerId: '',
                viewerName: '',
            }
        }),
    }
});
