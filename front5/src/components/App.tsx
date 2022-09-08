import React from "react";
import styles from "./App.module.css";
import {Router} from "./Router";
import {Header} from "./Header";

export function App() {
    return (
        <React.StrictMode>
            <div className={styles.app}>
                <Header/>
                <Router/>
            </div>
        </React.StrictMode>
    )
}