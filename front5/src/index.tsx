import {createRoot} from "react-dom/client";
import React from "react";
import {App} from "./components/App";
import {runForever} from "./store/onlineManager";

const root = document.getElementById('root');
if (root) {
    createRoot(root).render(<App/>);
}

runForever();
