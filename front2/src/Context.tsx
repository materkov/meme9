import React from "react";
import {Store} from "./Store";

export const GlobalStoreContext = React.createContext<Store>(new Store(() => {}));
