import {create} from "zustand";
import {AuthResp} from "../api/api";

export interface Globals {
    authToken: string;
    viewerId: string;
    viewerName: string;
    setAuth: (authResp: AuthResp) => void;
}

export const useGlobals = create<Globals>()(set => {
    let authToken = window.__prefetchApi.authToken;
    let viewerId = window.__prefetchApi.viewerId;
    let viewerName = window.__prefetchApi.viewerName;

    return {
        authToken: authToken,
        viewerId: viewerId,
        viewerName: viewerName,
        setAuth: (authResp: AuthResp) => {
            set({
                authToken: authResp.token,
                viewerId: authResp.userId,
                viewerName: authResp.userName,
            });
        }
    }
});
