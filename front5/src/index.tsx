import {createRoot} from "react-dom/client";
import React from "react";
import {Link} from "./components/Link";
import {Router} from "./components/Router";

const App = () => {
    return (
        <React.StrictMode>
            <Link href={"/"}>Feed</Link>
            <br/>
            <Router/>
        </React.StrictMode>
    );
}

const root = document.getElementById('root');
if (root) {
    createRoot(root).render(<App/>);
}
