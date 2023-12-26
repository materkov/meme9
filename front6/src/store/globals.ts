import {create} from "zustand";
import {AuthResp} from "../api/api";
import {tryGetPrefetch} from "../utils/prefetch";

export interface Globals {
    viewerId: string;
    viewerName: string;
    setAuth: (authResp: AuthResp) => void;
}

// TODO rewrite to React context
export const useGlobals = create<Globals>()(set => {
    const viewerId = tryGetPrefetch('viewerId') || "";
    const viewerName = tryGetPrefetch('viewerName') || "";

    return {
        viewerId: viewerId,
        viewerName: viewerName,
        setAuth: (authResp: AuthResp) => {
            set({
                viewerId: authResp.userId,
                viewerName: authResp.userName,
            });
        }
    }
});
