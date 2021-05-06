import React from "react";

export const GlobalContext = React.createContext((route: string) => {
})

export interface GlobalStore {
    togglePostLike(postId: string): void;
    followUser(userId: string): void;
    unfollowUser(userId: string): void;
}

export const GlobalStoreContext = React.createContext<GlobalStore | null>(null);
