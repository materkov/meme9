import React from "react";
import styles from "./App.module.css";
import {Router} from "./Router";
import {Header} from "./Header";
import {Provider} from "react-redux";
import {store} from "../store/store";

export function App() {
    return (
        <React.StrictMode>
            <Provider store={store}>
                <div className={styles.app}>
                    <Header/>
                    <Router/>
                </div>
            </Provider>
        </React.StrictMode>
    )
}