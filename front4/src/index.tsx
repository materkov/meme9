import React from "react";
import {createRoot} from "react-dom/client";
import {ResolveRoute} from "./routing";

const root = document.getElementById('root');
if (root) {
    createRoot(root).render(<ResolveRoute/>);
}
