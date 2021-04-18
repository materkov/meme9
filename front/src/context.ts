import React from "react";
import * as schema from "./schema/api";

export const NavigateContext = React.createContext<(url: string) => void>(() => {});
NavigateContext.displayName = "navigate";

export interface Store {
    Header: schema.HeaderRenderer;
    Page: schema.IndexRenderer | schema.LoginPageRenderer;
}
