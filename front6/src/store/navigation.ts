import {create} from "zustand";

export interface Navigation {
    url: string;
    setURL: (url: string) => void;
}

export const useNavigation = create<Navigation>()(set => ({
    url: window.location.pathname,
    setURL: (url: string) => set(() => {
        window.history.pushState(null, '', url);
        return {url: url}
    }),
}));
