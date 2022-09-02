import {createRoot} from "react-dom/client";
import React from "react";
import {Link} from "./components/Link";
import {Router} from "./components/Router";
import {QueryClientProvider} from "@tanstack/react-query";
import {queryClient} from "./store/store";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

const App = () => {
    return (
        <React.StrictMode>
            <QueryClientProvider client={queryClient}>
                <Link href={"/"}>Feed</Link>
                <br/>
                <Router/>
                <ReactQueryDevtools/>
            </QueryClientProvider>
        </React.StrictMode>
    );
}

const root = document.getElementById('root');
if (root) {
    createRoot(root).render(<App/>);
}
