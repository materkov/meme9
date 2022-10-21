import React from "react";
import styles from "./App.module.css";
import {Router} from "./Router";
import {Header} from "./Header";
import {queryClient} from "../store/fetcher";
import {QueryClientProvider} from "@tanstack/react-query";

export function App() {
    return (
        <React.StrictMode>
            <QueryClientProvider client={queryClient}>
                <div className={styles.app}>
                    <Header/>
                    <Router/>
                </div>
            </QueryClientProvider>
        </React.StrictMode>
    )
}