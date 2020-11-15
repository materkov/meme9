import React from "react";

export const NavigateContext = React.createContext<(url: string) => void>(() => {});
