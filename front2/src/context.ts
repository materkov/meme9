import React from "react";

export const UrlContext = React.createContext({
    url: '/',
    navigate: (url: string) => {
    },
});