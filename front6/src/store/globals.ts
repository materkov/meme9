import {create} from "zustand";

export interface Globals {
    viewerId: string;
}

export const useGlobals = create<Globals>()(set => ({
    viewerId: localStorage.getItem('viewerId') || "",
    setViewerId: (userId: string) => set(() => {
        localStorage.setItem('viewerId', userId);
        return {viewerId: userId};
    })
}));
